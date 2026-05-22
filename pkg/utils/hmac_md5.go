package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
)

// HmacMd5 生成十六进制编码的 HMAC-MD5 摘要。
func HmacMd5(key, data string) string {
	hash := hmac.New(md5.New, []byte(key))
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
