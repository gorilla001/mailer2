package utils

import (
	iconv "gopkg.in/iconv.v1"
)

// GBK2UTF8 convert gbk stirng into utf8
// TODO ensure input `gbk` is indeed gbk encoded
func GBK2UTF8(gbk string) string {
	conv, err := iconv.Open("utf-8", "gbk") // gbk -> utf-8
	if err != nil {
		return ""
	}
	defer conv.Close()
	return conv.ConvString(gbk)
}
