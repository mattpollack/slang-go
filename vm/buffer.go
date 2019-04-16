package vm

import (
	"fmt"
	"unsafe"
)

type ByteBuffer struct {
	bytes []byte
}

func NewByteBuffer() *ByteBuffer {
	return &ByteBuffer{
		[]byte{},
	}
}

func (b *ByteBuffer) Extend(i int) {
	b.bytes = append(b.bytes, make([]byte, i)...)
}

func (b *ByteBuffer) Len() int {
	return len(b.bytes)
}

func (b *ByteBuffer) Get8(i int) uint8 {
	if i+1 > len(b.bytes) {
		return 0
	}

	return *(*uint8)(unsafe.Pointer(&b.bytes[i]))
}

func (b *ByteBuffer) Get16(i int) uint16 {
	if i+2 > len(b.bytes) {
		return 0
	}

	return *(*uint16)(unsafe.Pointer(&b.bytes[i]))
}

func (b *ByteBuffer) Get32(i int) uint32 {
	if i+4 > len(b.bytes) {
		return 0
	}

	return *(*uint32)(unsafe.Pointer(&b.bytes[i]))
}

func (b *ByteBuffer) Get64(i int) uint64 {
	if i+8 > len(b.bytes) {
		return 0
	}

	return *(*uint64)(unsafe.Pointer(&b.bytes[i]))
}

func (b *ByteBuffer) Set8(i int, value uint8) {
	if i+1 >= len(b.bytes) {
		b.Extend(i + 1 - len(b.bytes))
	}

	*(*uint8)(unsafe.Pointer(&b.bytes[i])) = value
}

func (b *ByteBuffer) Set16(i int, value uint16) {
	if i+2 >= len(b.bytes) {
		b.Extend(i + 2 - len(b.bytes))
	}

	*(*uint16)(unsafe.Pointer(&b.bytes[i])) = value
}

func (b *ByteBuffer) Set32(i int, value uint32) {
	if i+4 >= len(b.bytes) {
		b.Extend(i + 4 - len(b.bytes))
	}

	*(*uint32)(unsafe.Pointer(&b.bytes[i])) = value
}

func (b *ByteBuffer) Set64(i int, value uint64) {
	if i+8 >= len(b.bytes) {
		b.Extend(i + 16 - len(b.bytes))
	}

	*(*uint64)(unsafe.Pointer(&b.bytes[i])) = value
}

func (b *ByteBuffer) Print() {
	fmt.Println(b.bytes)
}
