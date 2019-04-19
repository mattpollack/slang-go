package vm

import (
	"encoding/binary"
	"fmt"
)

const (
	_ = iota

	VM_RUNNING
	VM_DONE
	VM_ERROR
)

type VM struct {
	Prog      *ByteBuffer
	stack     *ByteBuffer
	callStack []int

	pp int

	Status int
	Err    error
}

func NewVM(prog *ByteBuffer) *VM {
	return &VM{
		prog,
		NewByteBuffer(),
		[]int{},
		0,
		VM_RUNNING,
		nil,
	}
}

func (v *VM) Pop() Instruction {
	return nil
}

func (v *VM) NextOP() OP {
	op := OP{}
	op.Decode(v.Prog.Slice(v.pp, -1))
	v.pp += op.Size()

	return op
}

func (v *VM) Error(err error) int {
	v.Status = VM_ERROR
	v.Err = err

	return v.Status
}

func (v *VM) Step() int {
	if v.pp >= v.Prog.Len() {
		return v.Error(fmt.Errorf("Program pointer %d is out of bounds", v.pp))
	}

	switch v.NextOP().Kind {
	case INS_PUSH:
		arg := &Data{}
		arg.Decode(v.Prog.Slice(v.pp, -1))

		// Use a bytebuffer to extract bytes
		tempBuffer := NewByteBuffer()
		arg.Emit(tempBuffer)

		v.stack.Push(tempBuffer.Bytes())

		v.pp += arg.Size()

	case INS_DUP:
		arg := &Data{}
		arg.Decode(v.stack.Slice(0, -1))

		// Use a bytebuffer to extract bytes
		tempBuffer := NewByteBuffer()
		arg.Emit(tempBuffer)
		v.stack.Push(tempBuffer.Bytes())

	case INS_POP:
		arg := &Data{}
		arg.Decode(v.stack)
		v.stack = v.stack.Slice(arg.Size(), -1)

	case INS_EXIT:
		v.Status = VM_DONE

	case INS_RETURN:
		v.pp = v.callStack[len(v.callStack)-1]
		v.callStack = v.callStack[:len(v.callStack)-1]

	case INS_CALL:
		arg := &Data{}
		arg.Decode(v.Prog.Slice(v.pp, -1))
		v.pp += arg.Size()

		// Use a bytebuffer to extract address
		tempBuffer := NewByteBuffer()
		tempBuffer.Push(arg.Value)

		// Set return address
		v.callStack = append(v.callStack, v.pp)

		// Set new program pointer
		v.pp = int(tempBuffer.Get32(0))

	case INS_ADD:
		arg0 := &Data{}
		arg0.Decode(v.stack.Slice(0, -1))
		v.stack = v.stack.Slice(arg0.Size(), -1)

		arg1 := &Data{}
		arg1.Decode(v.stack.Slice(0, -1))
		v.stack = v.stack.Slice(arg1.Size(), -1)

		v0 := int32(binary.LittleEndian.Uint32(arg0.Value))
		v1 := int32(binary.LittleEndian.Uint32(arg1.Value))

		res := NewInt32(v0 + v1)
		tempBuffer := NewByteBuffer()
		res.Emit(tempBuffer)
		v.stack.Push(tempBuffer.Bytes())

	case INS_JNE:
		arg0 := &Data{}
		arg0.Decode(v.stack)
		v.stack = v.stack.Slice(arg0.Size(), -1)

		arg1 := &Data{}
		arg1.Decode(v.stack)
		v.stack = v.stack.Slice(arg0.Size(), -1)

		addr := &Data{}
		addr.Decode(v.Prog.Slice(v.pp, -1))
		v.pp += addr.Size()

		// Jump if arguments aren't equal
		if !arg0.Equals(arg1) {
			// Use a bytebuffer to extract address
			tempBuffer := NewByteBuffer()
			tempBuffer.Push(addr.Value)

			// Set new program pointer
			v.pp = int(tempBuffer.Get32(0))
		}

	default:
		return v.Error(fmt.Errorf("Instruction '%s' doesn't have an implementation", INS_NAME[int(v.Prog.Get16(v.pp))]))
	}

	return v.Status
}

func (v *VM) PrintStack() {
	for i := v.stack.Len() - 1; i >= 0; i-- {
		fmt.Printf("| %d\n", v.stack.Get8(i))
	}

	fmt.Println(v.callStack)
}
