package cache

type nop struct{}

func newNop() *nop {
	return &nop{}
}

func (nop) Get(ket int64) (value interface{}, exists bool) { return nil, false }
func (nop) Remove(key int64)                               {}
func (nop) Add(key int64, value interface{})               {}
