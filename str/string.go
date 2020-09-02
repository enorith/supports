package str

import (
	"github.com/enorith/supports/byt"
	"strings"
	"unsafe"
)

func Contains(haystack string, needles ...string) bool {
	for _, n := range needles {
		if strings.Contains(haystack, n) {
			return true
		}
	}

	return false
}

func StartWith(haystack string, needle string) bool {
	return strings.HasPrefix(haystack, needle)
}

func EndWith(haystack string, needle string) bool {
	return strings.HasSuffix(haystack, needle)
}

func Duplicate(str string, times int) []string {
	var value []string
	for i := 0; i < times; i++ {
		value = append(value, str)
	}
	return value
}

func RandString(n int) string {
	b := byt.RandBytes(n)

	return *(*string)(unsafe.Pointer(&b))
}
