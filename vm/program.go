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
	INS_JMP

	INS_DUP
	INS_POP

	INS_PUSH
	INS_CALL
	INS_CALL_ENV
	INS_RETURN
	INS_EXIT

	T_INT_32
	T_STRING
	T_ADDRESS
	T_ENVIRONMENT
)

var INS_NAME map[int]string = map[int]string{
	INS_ADD: "add",
	INS_SUB: "sub",
	INS_MUL: "mul",
	INS_DIV: "div",

	INS_JNE: "jne",
	INS_JMP: "jmp",

	INS_DUP: "dup",
	INS_POP: "pop",

	INS_PUSH:     "push",
	INS_CALL:     "call",
	INS_CALL_ENV: "call_env",
	INS_RETURN:   "return",
	INS_EXIT:     "exit",

	T_INT_32:      "int32",
	T_STRING:      "string",
	T_ADDRESS:     "address",
	T_ENVIRONMENT: "environment",
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
