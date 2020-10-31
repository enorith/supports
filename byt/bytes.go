package byt

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"
)

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
const letterSeeds = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-=_+';./?><`~!@#$%^&*()"

func Contains(haystack []byte, needles ...[]byte) bool {
	for _, n := range needles {
		if bytes.Contains(haystack, n) {
			return true
		}
	}

	return false
}

//StartWith alias of bytes.HasPrefix
func StartWith(haystack []byte, needle []byte) bool {
	return bytes.HasPrefix(haystack, needle)
}

//EndWith alias of bytes.HasSuffix
func EndWith(haystack []byte, needle []byte) bool {
	return bytes.HasSuffix(haystack, needle)
}

func Duplicate(b []byte, times int) [][]byte {
	var value [][]byte
	for i := 0; i < times; i++ {
		value = append(value, b)
	}
	return value
}

var src = rand.NewSource(time.Now().UnixNano())

func RandBytes(n int) []byte {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterSeeds) {
			b[i] = letterSeeds[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}

func ToString(b []byte) string {
	return string(b)
}

func ToUint64(b []byte) (uint64, error) {
	str := ToString(b)

	return strconv.ParseUint(str, 10, 64)
}

func ToInt64(b []byte) (int64, error) {
	str := ToString(b)

	return strconv.ParseInt(str, 10, 64)
}

func ToBool(b []byte) (bool, error) {
	str := ToString(b)

	return strconv.ParseBool(str)
}

func ToFloat64(b []byte) (float64, error) {
	str := ToString(b)

	return strconv.ParseFloat(str, 64)
}
