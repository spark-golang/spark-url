package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"errors"
	"fmt"
	"strconv"
	"strings"

	hashids "github.com/speps/go-hashids"
)

const hashSalt = "jjKJIHLfShuJe9Vg"

// Sha1 sha1加密
func Sha1(str string) string {
	hash := sha1.New()
	hash.Write([]byte(str))
	bs := hash.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// HashDecode hashid解密
func HashDecode(value string) (string, error) {
	if len(value) != 12 {
		return "", errors.New("uid 不合法")
	}
	hd := hashids.NewData()
	hd.Salt = hashSalt
	hd.MinLength = 12
	h, _ := hashids.NewWithData(hd)

	d, err := h.DecodeWithError(value)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(d[0]), nil
}

// Md5 加密
func Md5(value string) string {
	return strings.ToLower(fmt.Sprintf("%x", md5.Sum([]byte(value))))
}

// 转换参数UID为整型
func TranslateUidParam(uidStr string) (uint32, error) {
	uidTmp, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		uidStr, err = HashDecode(uidStr)
		if err != nil {
			return 0, errors.New("用户ID错误")
		}

		uidTmp, err = strconv.ParseUint(uidStr, 10, 64)
		if err != nil {
			return 0, errors.New("用户ID错误")
		}
	}

	return uint32(uidTmp), nil
}
