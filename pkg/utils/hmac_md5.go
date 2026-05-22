package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
)

func HmacMd5(key, data string) string {
	hash := hmac.New(md5.New, []byte(key))
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
