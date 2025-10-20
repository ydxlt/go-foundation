package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

// Sha1Bytes 返回原始的 SHA1 二进制哈希
func Sha1Bytes(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}

// Sha1Hex 返回十六进制字符串
func Sha1Hex(data []byte) string {
	return hex.EncodeToString(Sha1Bytes(data))
}

// Sha1String 计算输入字符串的 SHA1 哈希并返回 hex 编码后的字符串
func Sha1String(s string) string {
	return Sha1Hex([]byte(s))
}
