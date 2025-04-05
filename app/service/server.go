package main

import (
	"fmt"
	"net"
	"study/app/service/pb"
	"study/xnet"
	"sync"
	"time"
)

type User struct {
	net.Conn
	userId   int
	roomId   int
	lastTime int64
}

type Users struct {
	m map[int]*User
	l sync.RWMutex
}

func ser() error {
	lister, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}

	for {
		conn, err := lister.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			}
			return err
		}

		go func(conn net.Conn) {

			var userConn = &User{Conn: conn}

			defer userConn.Close()

			for {
				_ = userConn.SetReadDeadline(time.Now().Add(time.Minute * 5))
				p, err := xnet.DecodePack(userConn)
				if err != nil {
					return
				}

				if err = dispatchPack(userConn, p); err != nil {
					return
				}
			}
		}(conn)
	}
}

func dispatchPack(user User, p *xnet.Pack) error {
	if user.userId == 0 && uint32(pb.Cmd_login) != p.Cmd {
		return fmt.Errorf("dodododo")
	}

	switch pb.Cmd(p.Cmd) {
	case pb.Cmd_login:
	case pb.Cmd_join:
	default:
		return fmt.Errorf("cmd %d not found", p.Cmd)

	}
	return nil
}

func main() {

}
