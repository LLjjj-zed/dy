package middleware

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/giorgisio/goav/avcodec"
	"github.com/giorgisio/goav/avformat"
	"github.com/giorgisio/goav/avutil"
	"github.com/giorgisio/goav/swscale"
	"image"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

// SaveFrame 将单个帧作为 PPM 文件写入磁盘
func SaveFrame(frame *avutil.Frame, width, height int, filename string) {
	// 打开文件
	var builder strings.Builder
	builder.WriteString("./tmp/")
	builder.WriteString(filename)
	builder.WriteString(".ppm")
	fileName := builder.String()
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("Error Reading")
	}
	defer file.Close()
	// 写入标头
	header := fmt.Sprintf("P6\n%d %d\n255\n", width, height)
	file.Write([]byte(header))
	// 写入像素数据
	for y := 0; y < height; y++ {
		data0 := avutil.Data(frame)[0]
		buf := make([]byte, width*3)
		startPos := uintptr(unsafe.Pointer(data0)) + uintptr(y)*uintptr(avutil.Linesize(frame)[0])
		for i := 0; i < width*3; i++ {
			element := *(*uint8)(unsafe.Pointer(startPos + uintptr(i)))
			buf[i] = element
		}
		file.Write(buf)
	}
}

// GetFirstStream 获取视频的第一帧
func GetFirstStream(filename string) error {
	//打开视频文件
	pFormatContext := avformat.AvformatAllocContext()
	if avformat.AvformatOpenInput(&pFormatContext, filename, nil, nil) != 0 {
		return errors.New("无法打开此文件")
	}
	//检索流信息
	if pFormatContext.AvformatFindStreamInfo(nil) < 0 {
		return errors.New("找不到流信息")
	}
	//将有关文件的信息转储到标准错误上
	pFormatContext.AvDumpFormat(0, filename, 0)
	//查找第一个视频流
	for i := 0; i < int(pFormatContext.NbStreams()); i++ {
		switch pFormatContext.Streams()[i].CodecParameters().AvCodecGetType() {
		case avformat.AVMEDIA_TYPE_VIDEO:
			// 获取指向视频流的编解码器上下文的指针
			pCodecCtxOrig := pFormatContext.Streams()[i].Codec()
			// 查找视频流的解码器
			pCodec := avcodec.AvcodecFindDecoder(avcodec.CodecId(pCodecCtxOrig.GetCodecId()))
			if pCodec == nil {
				return errors.New("不支持的编解码器！")
			}
			// 复制上下文
			pCodecCtx := pCodec.AvcodecAllocContext3()
			if pCodecCtx.AvcodecCopyContext((*avcodec.Context)(unsafe.Pointer(pCodecCtxOrig))) != 0 {
				return errors.New("无法复制编解码器上下文")
			}
			// 打开编解码器
			if pCodecCtx.AvcodecOpen2(pCodec, nil) < 0 {
				return errors.New("无法打开编解码器")
			}
			// 分配视频帧
			pFrame := avutil.AvFrameAlloc()
			// 分配 AVFrame 结构
			pFrameRGB := avutil.AvFrameAlloc()
			if pFrameRGB == nil {
				return errors.New("无法分配 RGB 帧")
			}
			// 确定所需的缓冲区大小并分配缓冲区
			numBytes := uintptr(avcodec.AvpictureGetSize(avcodec.AV_PIX_FMT_RGB24, pCodecCtx.Width(),
				pCodecCtx.Height()))
			buffer := avutil.AvMalloc(numBytes)
			// 将缓冲区的适当部分分配给 pFrameRGB 中的图像平面
			// 请注意，pFrameRGB 是 AVFrame，但 AVFrame 是 AVPicture 的超集。
			avp := (*avcodec.Picture)(unsafe.Pointer(pFrameRGB))
			avp.AvpictureFill((*uint8)(buffer), avcodec.AV_PIX_FMT_RGB24, pCodecCtx.Width(), pCodecCtx.Height())
			// 初始化 SWS 上下文以进行软件扩展
			swsCtx := swscale.SwsGetcontext(
				pCodecCtx.Width(),
				pCodecCtx.Height(),
				(swscale.PixelFormat)(pCodecCtx.PixFmt()),
				pCodecCtx.Width(),
				pCodecCtx.Height(),
				avcodec.AV_PIX_FMT_RGB24,
				avcodec.SWS_BILINEAR,
				nil,
				nil,
				nil,
			)
			// 读取帧并将第一帧保存到磁盘
			packet := avcodec.AvPacketAlloc()
			for pFormatContext.AvReadFrame(packet) >= 0 {
				// 这是来自视频流的数据包吗？
				if packet.StreamIndex() == i {
					// 解码视频帧
					response := pCodecCtx.AvcodecSendPacket(packet)
					if response < 0 {
						fmt.Printf("将数据包发送到解码器时出错: %s\n", avutil.ErrorFromCode(response))
					}
					for response >= 0 {
						response = pCodecCtx.AvcodecReceiveFrame((*avcodec.Frame)(unsafe.Pointer(pFrame)))
						if response == avutil.AvErrorEAGAIN || response == avutil.AvErrorEOF {
							break
						} else if response < 0 {
							fmt.Printf("从解码器接收帧时出错: %s\n", avutil.ErrorFromCode(response))
							return errors.New("从解码器接收帧时出错")
						}
						// 将图像从其本机格式转换为 RGB
						swscale.SwsScale2(swsCtx, avutil.Data(pFrame),
							avutil.Linesize(pFrame), 0, pCodecCtx.Height(),
							avutil.Data(pFrameRGB), avutil.Linesize(pFrameRGB))
						// 将帧保存到磁盘
						SaveFrame(pFrameRGB, pCodecCtx.Width(), pCodecCtx.Height(), filename)
					}
				}
				// 释放由av_read_frame分配的数据包
				packet.AvFreePacket()
			}
			// 释放 RGB 图像
			avutil.AvFree(buffer)
			avutil.AvFrameFree(pFrameRGB)
			// 释放 YUV 框架
			avutil.AvFrameFree(pFrame)
			// 关闭编解码器
			pCodecCtx.AvcodecClose()
			(*avcodec.Context)(unsafe.Pointer(pCodecCtxOrig)).AvcodecClose()
			// 关闭视频文件
			pFormatContext.AvformatCloseInput()
			// 保存第一个视频片段的帧后停止
			break
		default:
			return errors.New("找不到视频流")
		}
	}
	return nil
}

// 用于转换图像值的实用程序解析功能
func parseByte(strValue string) byte {
	var value uint64
	var err error

	if value, err = strconv.ParseUint(strValue, 10, 8); err != nil {
		fmt.Printf("Failed to convert %s: %s\n", strValue, err.Error())
		return byte(0)
	}
	return byte(value)
}

func ppm2png(filename string) error {
	var builder strings.Builder
	builder.WriteString("./tmp/")
	builder.WriteString(filename)
	builder.WriteString(".ppm")
	sourceFilename := builder.String()
	builder.Reset()
	builder.WriteString("./tmp/")
	builder.WriteString(filename)
	builder.WriteString(".png")
	destFilename := builder.String()
	defer os.Remove(sourceFilename)
	// 尝试获取文件
	readFile, err := os.Open(sourceFilename)
	if err != nil {
		fmt.Printf("Could not open file %s: %s\n", sourceFilename, err.Error())
		return err
	}
	defer readFile.Close()
	// 现在开始处理它
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	// 验证文件头
	const FileHeader = "P3"
	if fileScanner.Scan() && fileScanner.Text() != FileHeader {
		fmt.Printf("Failed to read %s, bad header\n", sourceFilename)
		return err
	}
	var values []string
	var width, height int
	if fileScanner.Scan() {
		line := fileScanner.Text()
		var value int64
		values = strings.Split(line, " ")
		if len(values) != 2 {
			fmt.Printf("Failed to read %s bad dimensions ->%d\n", sourceFilename, len(values))
			return err
		}
		if value, err = strconv.ParseInt(values[0], 10, 32); err != nil {
			fmt.Printf("Failed to read %s, bad dimensions: %s\n", sourceFilename, err.Error())
			return err
		}
		width = int(value)
		if value, err = strconv.ParseInt(values[0], 10, 32); err != nil {
			fmt.Printf("Failed to read %s, bad dimensions: %s\n", sourceFilename, err.Error())
			return err
		}
		height = int(value)
	}
	// 验证 BPP 值，即使不使用它们。
	if fileScanner.Scan() {
		line := fileScanner.Text()
		if _, err := strconv.ParseInt(line, 10, 32); err != nil {
			fmt.Printf("Failed to read %s, bad colour size: %s\n", sourceFilename, err.Error())
			return err
		}
	}
	// 准备数据缓冲区
	const BPP = 4
	imageData := image.NewRGBA(image.Rect(0, 0, width, height))
	row := 0
	for fileScanner.Scan() {
		line := fileScanner.Text()
		values := strings.Split(line, " ")
		if len(values) != 3 {
			fmt.Println(values)
			fmt.Printf("Failed to read %s at line %d, bad colour size\n", sourceFilename, row)
			return err
		}
		imageData.Pix[row] = parseByte(values[0])
		imageData.Pix[row+1] = parseByte(values[1])
		imageData.Pix[row+2] = parseByte(values[2])
		imageData.Pix[row+3] = 255
		// 前进到下一个像素组
		row += BPP
	}
	outputFile, err := os.Create(destFilename)
	if err != nil {
		fmt.Printf("Could not create file %s: %s\n", destFilename, err.Error())
		return err
	}
	defer outputFile.Close()
	png.Encode(outputFile, imageData)
	fmt.Printf("Processed %s into %s\n", sourceFilename, destFilename)
	return nil
}
