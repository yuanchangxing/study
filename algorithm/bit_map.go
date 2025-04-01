package algorithm

var bc = [64]uint64{}

func init() {
	bc[0] = 1
	for i := 1; i < 64; i++ {
		bc[i] = bc[i-1] << 1
	}
}

type bitMap struct {
	s []uint64
}

func (b *bitMap) Set(i int) {
	var index = i / 64
	for len(b.s) < index+1 {
		b.s = append(b.s, 0)
	}
	b.s[index] |= b._w(i)
}

func (b *bitMap) Clear(i int) {
	var index = i / 64
	b.s[index] &= ^b._w(i)
}

func (b *bitMap) _w(i int) uint64 {
	i = i % 64
	return bc[i]
}

func (b *bitMap) In(i int) bool {
	var index = i / 64
	return b.s[index]&b._w(i) > 0
}
