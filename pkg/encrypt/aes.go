package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"time"
)

// 使用用户密钥加密消息
func AseEncryptMessage(key []byte, message string) (string, error) {
	// 创建 AES 块加密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 将消息转为字节数组
	plaintext := []byte(message)

	// 创建一个随机的 IV (初始化向量)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	// 填充随机数据到 IV
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// 使用 CBC 模式加密
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// 返回加密后的消息（以十六进制编码）
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// 使用用户密钥解密消息
func AseDecryptMessage(key []byte, encryptedMessage string) (string, error) {
	// 将十六进制编码的加密消息转为字节数组
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedMessage)
	if err != nil {
		return "", err
	}

	// 创建 AES 块解密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 取出 IV
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// 使用 CBC 模式解密
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	// 返回解密后的明文
	return string(ciphertext), nil
}

func AesGenRandKey() string {
	types := []int{16, 24, 32}
	t := types[time.Now().Unix()%3]

	key, err := uuid.GenerateRandomBytes(t)
	if err != nil {
		logx.Error(err)
		return ""
	}

	return string(key)
}
