package xnet

import (
	"bytes"
	"log"
	"testing"
)

func TestEncodePack(t *testing.T) {
	var data = NewPack(1, []byte("hello world"))
	bs := EncodePack(data)
	w := bytes.NewBuffer(bs)
	p, err := DecodePack(w)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(string(p.Body))
}
