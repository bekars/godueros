package godueros

import (
	"math"
	"encoding/binary"
)

const (
	Sample_Rate 		int32 = 16000
	// 填满一个buffer调用一次回调
	Frame_Per_Buffer 	int   = 256
	Channels_Num		int32 = 1
)

var (
	access_token = "26.687423ae6cd02c090ef92af1b64be322.2592000.1514813398.2300547068-10218111"
	device_id = "wh9foSD4KnZAa8kyG1l62SdwkNHlVfuH"
)


func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func float32ToByte(f float32) []byte {
	bits := math.Float32bits(f)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func byteToFloat32(b []byte) float32 {
	bits := binary.LittleEndian.Uint32(b)
	return math.Float32frombits(bits)
}

func int16ToByte(n int16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(n))
	return b
}

func byteToInt16(b []byte) int16 {
	bits := binary.LittleEndian.Uint16(b)
	return int16(bits)
}

