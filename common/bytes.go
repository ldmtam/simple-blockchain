package common

import (
	"encoding/binary"
)

// Uint64 encodes []byte.
func Uint64(data []byte) uint64 {
	return binary.BigEndian.Uint64(data)
}

// FromUint64 decodes uint64 value.
func FromUint64(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

// Int64 encodes int64.
func Int64(data []byte) int64 {
	return int64(binary.BigEndian.Uint64(data))
}

// FromInt64 decodes int64.
func FromInt64(v int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// Uint32 encodes uint32
func Uint32(data []byte) uint32 {
	return binary.BigEndian.Uint32(data)
}

// FromUint32 decodes uint32
func FromUint32(v uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, v)
	return b
}
