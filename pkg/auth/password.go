package auth

import (
	"strings"
	"xhs-srv-main/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func VerifyPassword(password, hash string, salt string) bool {
	// 判断密码是否是 bcrypt 格式
	if strings.HasPrefix(hash, "$2a$") || strings.HasPrefix(hash, "$2b$") {
		// 使用 bcrypt 校验密码
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		if err != nil {
			return false
		}
	} else {
		// 使用 MD5 方式校验密码
		encryptedPassword := EncryptKdsPassword(password, salt)
		if encryptedPassword != hash {
			return false
		}
	}

	return true
}

func EncryptKdsPassword(password, salt string) string {
	return utils.Md5Hash(utils.Md5Hash(utils.Md5Hash(password)) + salt + "kkkdddssslife")
}
