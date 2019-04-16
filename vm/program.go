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

	T_INT_32
)

var INS_NAME map[int]string = map[int]string{
	INS_ADD: "ADD",
	INS_SUB: "SUB",
	INS_MUL: "MUL",
	INS_DIV: "DIV",

	T_INT_32: "INT32",
}

type Program struct {
	Blocks []*BasicBlock
}

func (p *Program) Push(block *BasicBlock) {
	p.Blocks = append(p.Blocks, block)
}

func (p *Program) Print() {
	for _, b := range p.Blocks {
		b.Print()
		fmt.Println()
	}
}

type BasicBlock struct {
	Label string
	Body  []Instruction
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

type Instruction interface {
	Emit(ByteBuffer)
	Size() int
	Print()
}

// --------------------------------------------------------

type TODO struct {
	Thing string
}

func (in *TODO) Emit(buffer ByteBuffer) {}

func (in *TODO) Size() int {
	return 0
}

func (in *TODO) Print() {
	fmt.Printf("  TODO %s\n", in.Thing)
}

// --------------------------------------------------------
type OP struct {
	Kind uint16
}

func (in *OP) Emit(buffer ByteBuffer) {
	buffer.Set16(buffer.Len(), in.Kind)
}

func (in *OP) Size() int {
	return 16
}

func (in *OP) Print() {
	fmt.Printf("  %s\n", INS_NAME[int(in.Kind)])
}

// --------------------------------------------------------

func NewInt32(val int32) *Data {
	buff := []byte{0, 0, 0, 0}

	for i := 0; i < 32/8; i++ {
		buff[3-i] = byte(val >> uint(i*8))
	}

	return &Data{
		T_INT_32,
		buff,
	}
}

type Data struct {
	Kind  int
	Value []byte
}

func (in *Data) Emit(buffer ByteBuffer) {
	buffer.Set8(buffer.Len(), uint8(in.Kind))

	for _, b := range in.Value {
		buffer.Set8(buffer.Len(), uint8(b))
	}
}

func (in *Data) Size() int {
	return len(in.Value)*8 + 8
}

func (in *Data) Print() {
	fmt.Printf("  %s ", INS_NAME[int(in.Kind)])
	fmt.Println(in.Value)
}
