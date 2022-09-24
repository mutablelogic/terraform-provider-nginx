package util

import (
	"crypto/md5"
	"fmt"
)

func MD5Hash(data []byte) string {
	return fmt.Sprintf("%X", md5.Sum(data))
}
