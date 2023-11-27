package storage

type MemUnit string

const (
	Byte     MemUnit = "b"
	Kilobyte MemUnit = "Kb"
	Megabyte MemUnit = "Mb"
	Gigabyte MemUnit = "Gb"
)

type Size struct {
	Qty  uint
	Unit MemUnit
}

type Handle struct {
	Header  interface{}
	Content []byte
	Length  int
}

// TODO deallocate
type Allocator interface {
	Request(Size) (*Handle, error)
	GrowTo(Size, *Handle) error
}
