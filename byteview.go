package GOcache

type Byteview struct {
	b []byte
}

func (c Byteview) Len() int {
	return len(c.b)
}

func (c Byteview) ByteSlice() []byte {
	return Clonebyte(c.b)
}

func (c Byteview) String() string {
	return string(c.b)
}

func Clonebyte(b []byte) []byte {
	nby := make([]byte, len(b))
	copy(nby, b)
	return nby
}
