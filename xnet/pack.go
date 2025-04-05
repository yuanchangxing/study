package xnet

import (
	"encoding/binary"
	"fmt"
	"io"
)

const PackMagicNum = 0x01ab

type Pack struct {
	MagicNum uint32
	Cmd      uint32
	Length   uint32 // body长度
	Body     []byte
}

func NewPack(cmd uint32, body []byte) *Pack {
	var pack = new(Pack)
	pack.MagicNum = PackMagicNum
	pack.Cmd = cmd
	pack.Length = uint32(len(body))
	pack.Body = body
	return pack
}

func EncodePack(pack *Pack) []byte {
	var head = [4 * 3]byte{}
	binary.BigEndian.PutUint32(head[0:4], uint32(pack.MagicNum))
	binary.BigEndian.PutUint32(head[4:8], uint32(pack.Cmd))
	binary.BigEndian.PutUint32(head[8:12], uint32(pack.Length))

	var bs = make([]byte, 0, 4*3+len(pack.Body))
	bs = append(bs, head[:]...)
	bs = append(bs, pack.Body...)
	return bs
}

func DecodePack(reader io.Reader) (*Pack, error) {
	var headBuff = make([]byte, 12)
	size, err := io.ReadFull(reader, headBuff)
	if err != nil {
		return nil, err
	}
	var p = &Pack{}
	p.MagicNum = binary.BigEndian.Uint32(headBuff[0:4])
	p.Cmd = binary.BigEndian.Uint32(headBuff[4:8])
	p.Length = binary.BigEndian.Uint32(headBuff[8:12])

	if p.MagicNum != PackMagicNum {
		return nil, fmt.Errorf("magic number error  %d %d %v", size, p.MagicNum, headBuff)
	}

	p.Body = make([]byte, p.Length)
	_, err = io.ReadFull(reader, p.Body)
	if err != nil {
		return nil, err
	}

	return p, nil
}
