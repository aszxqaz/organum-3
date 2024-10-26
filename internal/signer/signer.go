package signer

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func Sign(payload, secret string) string {
	hmac := hmac.New(sha256.New, []byte(secret))
	hmac.Write([]byte(payload))
	dataHmac := hmac.Sum(nil)
	hmacHex := hex.EncodeToString(dataHmac)
	return hmacHex
}

func Verify(payload, signature, secret string) bool {
	expectedSignature := Sign(payload, secret)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func Checksum(file []byte) string {
	fmt.Println(len(file))
	hashBytes := md5.Sum(file)
	return hex.EncodeToString(hashBytes[:])
}
