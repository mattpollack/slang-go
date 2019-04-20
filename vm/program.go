package vm

/*

BYTECODE SPECIFICATION:



*/

import (
	"fmt"
)

const (
	_ = iota
	INS_ADD
	INS_SUB
	INS_MUL
	INS_DIV

	INS_JNE

	INS_DUP
	INS_POP

	INS_PUSH
	INS_CALL
	INS_RETURN
	INS_EXIT

	T_INT_32
	T_STRING
	T_ADDRESS
)

var INS_NAME map[int]string = map[int]string{
	INS_ADD: "add",
	INS_SUB: "sub",
	INS_MUL: "mul",
	INS_DIV: "div",

	INS_JNE: "jne",

	INS_DUP: "dup",
	INS_POP: "pop",

	INS_PUSH:   "push",
	INS_CALL:   "call",
	INS_RETURN: "return",
	INS_EXIT:   "exit",

	T_INT_32:  "int32",
	T_STRING:  "string",
	T_ADDRESS: "address",
}

type Program struct {
	Blocks []*BasicBlock
}

func (p *Program) Push(block *BasicBlock) {
	p.Blocks = append(p.Blocks, block)
}

func (p *Program) Print() {
	for i, b := range p.Blocks {
		b.Print()

		if i != len(p.Blocks)-1 {
			fmt.Println()
		}
	}
}

func (p *Program) Size() int {
	size := 0

	for _, b := range p.Blocks {
		size += b.Size()
	}

	return size
}

func (p *Program) Render() *ByteBuffer {
	buffer := NewByteBuffer()
	lookup := map[string]uint32{}
	offset := uint32(0)

	// Build label lookup
	for _, block := range p.Blocks {
		lookup[block.Label] = offset
		offset += uint32(block.Size())
	}

	// Emit instruction bytecode
	for _, block := range p.Blocks {
		for _, ins := range block.Body {
			switch addr := ins.(type) {
			case *Address:
				NewAddress(lookup[addr.Label]).Emit(buffer)
			default:
				ins.Emit(buffer)
			}
		}
	}

	return buffer
}

type BasicBlock struct {
	Label string
	Body  []Instruction

	// for use by the program to render addresses
	offset int
}

func (b *BasicBlock) Push(i Instruction) {
	b.Body = append(b.Body, i)
}

func (b *BasicBlock) Print() {
	fmt.Printf("%s:\n", b.Label)

	for _, i := range b.Body {
		i.Print()
	}
}

func (b *BasicBlock) Size() int {
	size := 0

	for _, b := range b.Body {
		size += b.Size()
	}

	return size
}

type Instruction interface {
	Decode(*ByteBuffer)
	Emit(*ByteBuffer)
	Size() int
	Print()
}

// --------------------------------------------------------

type TODO struct {
	Thing string
}

func (in *TODO) Decode(*ByteBuffer) {}
func (in *TODO) Emit(*ByteBuffer)   {}

func (in *TODO) Size() int {
	return 0
}

func (in *TODO) Print() {
	fmt.Printf("  # TODO: %s\n", in.Thing)
}

// --------------------------------------------------------

type OP struct {
	Kind uint16
}

func (in *OP) Decode(buffer *ByteBuffer) {
	in.Kind = buffer.Get16(0)
}

func (in *OP) Emit(buffer *ByteBuffer) {
	buffer.Set16(buffer.Len(), in.Kind)
}

func (in *OP) Size() int {
	return 16 / 8
}

func (in *OP) Print() {
	fmt.Printf("  %s\n", INS_NAME[int(in.Kind)])
}

// --------------------------------------------------------

type Address struct {
	Label string
}

func (in *Address) Decode(buffer *ByteBuffer) {}
func (in *Address) Emit(buffer *ByteBuffer)   {}

func (in *Address) Size() int {
	return (16 + 32 + 32) / 8
}

func (in *Address) Print() {
	fmt.Printf("  &(%s)\n", in.Label)
}

// --------------------------------------------------------

func NewInt32(val int32) *Data {
	buff := []byte{0, 0, 0, 0}

	for i := 0; i < 32/8; i++ {
		buff[i] = byte(val >> uint(i*8))
	}

	return &Data{
		T_INT_32,
		buff,
	}
}

func NewString(str string) *Data {
	return &Data{
		T_STRING,
		[]byte(str),
	}
}

func NewAddress(val uint32) *Data {
	buff := []byte{0, 0, 0, 0}

	for i := 0; i < 32/8; i++ {
		buff[i] = byte(val >> uint(i*8))
	}

	return &Data{
		T_ADDRESS,
		buff,
	}
}

type Data struct {
	Kind  int
	Value []byte
}

func (in *Data) Decode(buffer *ByteBuffer) {
	in.Kind = int(buffer.Get16(0))

	len := int(buffer.Get32(2))

	in.Value = buffer.Slice(6, 6+len).Bytes()
}

func (a *Data) Equals(b *Data) bool {
	if len(a.Value) != len(b.Value) {
		return false
	}

	for i, _ := range a.Value {
		if a.Value[i] != b.Value[i] {
			return false
		}
	}

	return true
}

func (in *Data) Emit(buffer *ByteBuffer) {
	buffer.Set16(buffer.Len(), uint16(in.Kind))
	buffer.Set32(buffer.Len(), uint32(len(in.Value)))

	for _, b := range in.Value {
		buffer.Set8(buffer.Len(), uint8(b))
	}
}

func (in *Data) Size() int {
	return (16+32)/8 + len(in.Value)
}

func (in *Data) Print() {
	fmt.Printf("  %s (%d) ", INS_NAME[int(in.Kind)], len(in.Value)*8)
	fmt.Println(in.Value)
}
