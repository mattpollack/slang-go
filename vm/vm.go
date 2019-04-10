package vm

/*
import (
	"fmt"
)

const (
	_ = iota

	I_POP_32
	I_PUSH_32
	I_DUP_32

	I_ADD_32
	I_SUB_32
	I_MUL_32
	I_DIV_32

	I_GOTO_32 // 32 bit immediate address

	I_EXIT

	VM_RUNNING
	VM_ERROR
	VM_EXITED
)

var I_STR map[int]string = map[int]string{
	I_POP_32:  "POP_32",
	I_PUSH_32: "PUSH_32",
	I_DUP_32:  "DUP_32",

	I_ADD_32: "ADD_32",
	I_SUB_32: "SUB_32",
	I_MUL_32: "MUL_32",
	I_DIV_32: "DIV_32",

	I_GOTO_32: "GOTO_32",

	I_EXIT: "EXIT",
}

type VM struct {
	Prog  *ByteBuffer
	stack *ByteBuffer

	pp int
	sp int

	status int
	Err    error
}

func NewVM() *VM {
	return &VM{
		NewByteBuffer(),
		NewByteBuffer(),
		0,
		0,
		VM_RUNNING,
		nil,
	}
}

func (v *VM) SetInstruction(code uint8) {
	v.Prog.Set8(v.Prog.Len(), code)
}

func (v *VM) Error(err error) int {
	v.status = VM_ERROR
	v.Err = err

	return v.status
}

func (v *VM) Step() int {
	if v.pp >= v.Prog.Len() {
		return v.Error(fmt.Errorf("Program pointer %d is out of bounds", v.pp))
	}

	switch v.Prog.Get8(v.pp) {
	case I_ADD_32:
		a := int(v.stack.Get32(v.sp - 4))
		b := int(v.stack.Get32(v.sp - 8))

		v.stack.Set32(v.sp-8, uint32(a+b))
		v.sp -= 4
		v.pp += 1

	case I_PUSH_32:
		v.stack.Set32(v.sp, v.Prog.Get32(v.pp+1))
		v.pp += 5
		v.sp += 4

	case I_POP_32:
		v.sp -= 4
		v.pp += 1

	case I_DUP_32:
		v.stack.Set32(v.sp, v.stack.Get32(v.sp-4))
		v.sp += 4
		v.pp += 1

	case I_GOTO_32:
		v.pp = int(v.stack.Get32(v.sp - 4))
		v.sp -= 4

	case I_EXIT:
		v.status = VM_EXITED
		v.pp += 1

	default:
		return v.Error(fmt.Errorf("Instruction with code %d doesn't exist", v.Prog.Get16(v.pp)))
	}

	return v.status
}

func (v *VM) PrintProg() {
	fmt.Println(v.Prog)
}

func (v *VM) PrintStack() {
	for i := v.sp - 1; i >= 0; i-- {
		fmt.Printf("| %d\n", v.stack.Get8(i))
	}

	fmt.Println("-")
}
*/
