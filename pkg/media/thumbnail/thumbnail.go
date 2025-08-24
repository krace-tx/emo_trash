package thumbnail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
)

// Thumbnail 代表一个图像对象，包含图像的属性，如大小、质量等
type Thumbnail struct {
	Source       image.Image // 原图
	OutputWidth  uint        // 输出的宽度
	OutputHeight uint        // 输出的高度
	Quality      int         // 输出质量（0-100）
	Format       string      // 输出格式
}

// NewThumbnail 创建一个新的 Thumbnail 实例
func NewThumbnail(source image.Image) *Thumbnail {
	return &Thumbnail{
		Source:  source,
		Quality: 80,     // 默认质量
		Format:  "auto", // 默认格式
	}
}

// LoadFromFile 从文件加载图像
func LoadFromFile(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}

// LoadFromBytes 从字节数组加载图像
func LoadFromBytes(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	return img, err
}

// Resize 设置输出图像的宽度和高度
func (tb *Thumbnail) Resize(width, height uint) *Thumbnail {
	tb.OutputWidth = width
	tb.OutputHeight = height
	return tb
}

// SetQuality 设置图像压缩后的质量（0-100）
func (tb *Thumbnail) SetQuality(quality int) *Thumbnail {
	if quality < 1 || quality > 100 {
		quality = 80 // 默认值
	}
	tb.Quality = quality
	return tb
}

// Compression 压缩图像并返回字节流（压缩后的图像）
// 保留原图格式进行压缩
func (tb *Thumbnail) Compression(quality int) ([]byte, error) {
	tb.Quality = quality
	var buf bytes.Buffer

	// 获取图像格式类型
	imgType := tb.Source

	// 如果需要，先对图像进行缩放
	if tb.OutputWidth > 0 && tb.OutputHeight > 0 {
		tb.Source = resizeImage(tb.Source, int(tb.OutputWidth), int(tb.OutputHeight))
	}

	// 根据图像类型（格式）决定压缩方法
	switch imgType := imgType.(type) {
	case *image.RGBA, *image.NRGBA, *image.RGBA64:
		// JPG 格式处理
		err := jpeg.Encode(&buf, tb.Source, &jpeg.Options{Quality: tb.Quality})
		if err != nil {
			return nil, fmt.Errorf("转换为JPG格式失败: %v", err)
		}
	case *image.Paletted:
		// PNG 格式处理
		err := png.Encode(&buf, tb.Source)
		if err != nil {
			return nil, fmt.Errorf("转换为PNG格式失败: %v", err)
		}
	case *image.YCbCr:
		// JPG 格式处理
		err := jpeg.Encode(&buf, tb.Source, &jpeg.Options{Quality: tb.Quality})
		if err != nil {
			return nil, fmt.Errorf("转换为JPG格式失败: %v", err)
		}
	default:
		return nil, fmt.Errorf("不支持的图像格式: %T", imgType)
	}

	return buf.Bytes(), nil
}

// resizeImage 根据指定的宽度和高度缩放图像
func resizeImage(img image.Image, width, height int) image.Image {
	// 你可以在这里实现你自己的图像缩放逻辑（例如使用 `github.com/nfnt/resize` 库）
	return img // 占位符：实际需要实现图像缩放
}

// 将图片编码为 Base64
func EncodeImageToBase64(imagePath string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %v", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %v", err)
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

// 将图片 URL 内容编码为 Base64
func EncodeImageURLToBase64(imageURL string) (string, error) {
	// 发起 HTTP GET 请求获取图片内容
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image from URL: %v", err)
	}
	defer resp.Body.Close()

	// 检查 HTTP 响应状态
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch image: HTTP status %d", resp.StatusCode)
	}

	// 读取图片数据
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data from URL: %v", err)
	}

	// 编码为 Base64
	return base64.StdEncoding.EncodeToString(data), nil
}
