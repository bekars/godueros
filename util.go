package godueros

import (
	"math"
	"encoding/binary"
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

