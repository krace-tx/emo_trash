package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

// RSAEncryptor 封装了 RSA 加密解密功能
type RSAEncryptor struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// NewRSAEncryptor 生成一个新的 RSA 加密解密器
func NewRSAEncryptor(bits int) (*RSAEncryptor, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %v", err)
	}

	return &RSAEncryptor{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}, nil
}

// Encrypt 使用公钥进行加密
func (e *RSAEncryptor) Encrypt(plainText []byte) ([]byte, error) {
	cipherText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, e.PublicKey, plainText, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt: %v", err)
	}
	return cipherText, nil
}

// Decrypt 使用私钥进行解密
func (e *RSAEncryptor) Decrypt(cipherText []byte) ([]byte, error) {
	plainText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, e.PrivateKey, cipherText, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %v", err)
	}
	return plainText, nil
}

// ExportPrivateKeyToPEM 导出私钥为 PEM 格式
func (e *RSAEncryptor) ExportPrivateKeyToPEM() ([]byte, error) {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(e.PrivateKey)
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	return pem.EncodeToMemory(pemBlock), nil
}

// ExportPublicKeyToPEM 导出公钥为 PEM 格式
func (e *RSAEncryptor) ExportPublicKeyToPEM() ([]byte, error) {
	pubKeyBytes := x509.MarshalPKCS1PublicKey(e.PublicKey)
	pemBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyBytes,
	}
	return pem.EncodeToMemory(pemBlock), nil
}

// LoadPrivateKeyFromPEM 从 PEM 格式的私钥加载
func (e *RSAEncryptor) LoadPrivateKeyFromPEM(pemData []byte) error {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return errors.New("failed to parse PEM block containing the private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}

	e.PrivateKey = privateKey

	return nil
}

// LoadPublicKeyFromPEM 从 PEM 格式的公钥加载
func (e *RSAEncryptor) LoadPublicKeyFromPEM(pemData []byte) error {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return errors.New("failed to parse PEM block containing the public key")
	}

	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %v", err)
	}

	e.PublicKey = pubKey

	return nil
}
