package vm

/*

BYTECODE SPECIFICATION:



*/

const (
	_ = iota

	BC_POP
)

type Program struct {
	buffer ByteBuffer
}

type ByteCode interface {
	Type() uint64 // Some method for making sure these types are unique
	Emit() []byte
}
