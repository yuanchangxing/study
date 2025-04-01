package algorithm

import (
	"log"
	"testing"
)

func TestBitMap(t *testing.T) {
	var s = bitMap{}
	s.Set(64)
	log.Println(s.In(64))
	s.Clear(64)
	log.Println(s.In(64))
}
