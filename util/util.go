package util

import "unsafe"

// StringToBytes 将字符串转换成byte数组，注意返回的数组为只读，不可修改
func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesToString 将byte数组转换成字符串
func BytesToString(bs []byte) string {
	return unsafe.String(unsafe.SliceData(bs), len(bs))
}
