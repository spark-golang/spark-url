package utils

import (
	"bytes"
	"encoding/binary"
)

// Bytes2Int 字节转整形
func Bytes2Int(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int32
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return int(tmp)
}

// Int2Bytes 整形转换成字节
func Int2Bytes(n int) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}

// Byte2UInt64 byte -> uint64
func Byte2UInt64(buf []byte) uint64 {
	return binary.BigEndian.Uint64(buf)
}

// UInt64ToByte uint64 -> byte
func UInt64ToByte(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}
