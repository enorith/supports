package carbon

import "time"

func GetMillisecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetTimestamp() int64 {
	return time.Now().Unix()
}

func FormatNow(format string) string {
	return time.Now().Format(format)
}
