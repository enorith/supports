package str

import (
	"strings"
	"unsafe"

	"github.com/enorith/supports/byt"
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

func RandString(n int, seeds ...string) string {
	b := byt.RandBytes(n, seeds...)

	return *(*string)(unsafe.Pointer(&b))
}

func ReplaceVar(str string, vars map[string]string) string {
	for k, v := range vars {
		str = strings.ReplaceAll(str, ":"+k, v)
	}

	return str
}
