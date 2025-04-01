package algorithm

import (
	"log"
	"testing"
)

func TestBigAdd(t *testing.T) {
	log.Println(BigAdd("1", "1"))
	log.Println(BigAdd("9", "9"))
	log.Println(BigAdd("19", "9"))
	log.Println(BigAdd("19", "0"))
	log.Println(BigAdd("19", "91"))
	log.Println(BigAdd("19", "9x"))
}
