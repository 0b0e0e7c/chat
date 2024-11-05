package common

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

func Hashing(username, password string) string {
	hash := sha256.New()

	hash.Write([]byte(username + password))

	return hex.EncodeToString(hash.Sum(nil))
}

func HashingUserID(uid int64) string {
	hash := sha256.New()

	hash.Write([]byte(strconv.FormatInt(uid, 10)))

	return hex.EncodeToString(hash.Sum(nil))
}
