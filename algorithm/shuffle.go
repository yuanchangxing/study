package algorithm

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"os"
	"runtime"
)

func Init(seed int64) error {
	if dr != nil {
		return nil
	}
	dr = &rang{}
	if runtime.GOOS == "windows" {
		dr.rng = rand.New(rand.NewSource(seed))
		return nil
	}
	var err error
	dr.file, err = os.Open("/dev/urandom")
	if err != nil {
		dr = nil
	}
	return err
}

type rang struct {
	rng  *rand.Rand
	file *os.File
}

var dr *rang

func (dr *rang) Int64n(n int64) int64 {
	if n <= 0 {
		panic("invalid argument to Int63n")
	}

	if n&(n-1) == 0 { // n is power of two, can mask
		return dr.int64() & (n - 1)
	}

	max := int64((1 << 63) - 1 - (1<<63)%uint64(n))

	v := dr.int64()
	for v > max {
		v = dr.int64()
	}
	return v % n
}

func (dr *rang) int64() int64 {
	if dr.rng != nil {
		return dr.rng.Int63()
	}

	var n uint64
	bs := make([]byte, 8)
	dr.file.Read(bs)
	rd := bytes.NewReader(bs)
	binary.Read(rd, binary.LittleEndian, &n)
	return int64(n >> 1)
}

// Fisher-Yates shuffles
func Shuffle(list []int) []int {
	n := len(list)

	for i := n - 1; i > 0; i-- {
		j := dr.Int64n(int64(i + 1))
		list[i], list[j] = list[j], list[i]
	}

	return list
}

func Int64n(n int64) int64 {
	return dr.Int64n(n)
}
