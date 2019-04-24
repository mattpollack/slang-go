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

func (v *VM) PopData() Data {
	data := Data{}
	data.Decode(v.stack.Slice(0, -1))
	v.stack = v.stack.Slice(data.Size(), -1)

	return data
}

func (v *VM) PushData(data Data) {
	tempBuffer := NewByteBuffer()
	data.Emit(tempBuffer)
	v.stack.Push(tempBuffer.Bytes())
}

func (v *VM) NextOP() OP {
	op := OP{}
	op.Decode(v.Prog.Slice(v.pp, -1))
	v.pp += op.Size()

	return op
}

func (v *VM) NextData() Data {
	data := Data{}
	data.Decode(v.Prog.Slice(v.pp, -1))
	v.pp += data.Size()

	return data
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
		// Push immediate value
		v.PushData(v.NextData())

	case INS_DUP:
		// Decode the top of stack
		arg := Data{}
		arg.Decode(v.stack.Slice(0, -1))

		// Push decoded data
		v.PushData(arg)

	case INS_POP:
		// Pop and discard
		v.PopData()

	case INS_EXIT:
		v.Status = VM_DONE

	case INS_RETURN:
		// Set program pointer and pop callstack
		v.pp = v.callStack[len(v.callStack)-1]
		v.callStack = v.callStack[:len(v.callStack)-1]

	case INS_CALL:
		// Get address by converting
		addr := int(binary.LittleEndian.Uint32(v.NextData().Value))

		// Set return address
		v.callStack = append(v.callStack, v.pp)

		// Set new program pointer
		v.pp = addr

	case INS_ADD:
		// Pop arguments
		arg0 := v.PopData()
		arg1 := v.PopData()

		// Decode values
		v0 := int32(binary.LittleEndian.Uint32(arg0.Value))
		v1 := int32(binary.LittleEndian.Uint32(arg1.Value))

		// Generate and push new value
		v.PushData(*NewInt32(v0 + v1))

	case INS_JMP:
		v.pp = int(binary.LittleEndian.Uint32(v.NextData().Value))

	case INS_JNE:
		// Pop arguments
		arg0 := v.PopData()
		arg1 := v.PopData()
		addr := v.NextData()

		// Jump if arguments aren't equal
		if !arg0.Equals(&arg1) {
			// Set new program pointer
			v.pp = int(binary.LittleEndian.Uint32(addr.Value))
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
