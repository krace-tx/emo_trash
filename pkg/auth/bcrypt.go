package authx

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// 盐值长度（字节），推荐16-32字节
const saltLength = 16

// GenerateSalt 生成随机盐值（base64编码字符串）
// 返回盐值字符串和错误（如随机数生成失败）
func GenerateSalt() (string, error) {
	// 生成指定长度的随机字节
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", errors.New("盐值生成失败: " + err.Error())
	}
	// 编码为base64字符串（便于存储）
	return base64.URLEncoding.EncodeToString(salt), nil
}

// HashPassword 结合盐值生成密码哈希
// password: 原始密码
// salt: 盐值（GenerateSalt返回的字符串）
// 返回哈希后的密码字符串和错误
func HashPassword(password, salt string) (string, error) {
	if password == "" || salt == "" {
		return "", errors.New("密码和盐值不能为空")
	}
	// 拼接密码与盐值（格式：password + salt）
	passwordWithSalt := []byte(password + salt)
	// 使用bcrypt生成哈希（cost=10，默认值）
	hash, err := bcrypt.GenerateFromPassword(passwordWithSalt, bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("密码哈希失败: " + err.Error())
	}
	return string(hash), nil
}

// VerifyPassword 验证密码（结合盐值）
// password: 待验证的原始密码
// salt: 存储的盐值
// hashedPassword: 存储的哈希密码
// 返回验证结果（true为验证通过）和错误
func VerifyPassword(password, salt, hashedPassword string) (bool, error) {
	if password == "" || salt == "" || hashedPassword == "" {
		return false, errors.New("密码、盐值和哈希值不能为空")
	}
	// 拼接密码与盐值
	passwordWithSalt := []byte(password + salt)
	// 比对哈希值
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), passwordWithSalt)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil // 密码不匹配（无错误）
	} else if err != nil {
		return false, errors.New("密码验证失败: " + err.Error()) // 其他错误（如哈希格式错误）
	}
	return true, nil
}
