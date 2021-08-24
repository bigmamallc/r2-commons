package tcp

import (
	"github.com/phayes/freeport"
	"net"
	"time"
)

func Kill(conn *net.TCPConn) {
	_ = conn.SetLinger(0)
	_ = conn.Close()
}

func PokeAddr(addr *net.TCPAddr) bool {
	conn, err := net.DialTCP("tcp4", nil, addr)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func Await(addr *net.TCPAddr, deadline time.Time) bool {
	for time.Now().Before(deadline) {
		if PokeAddr(addr) {
			return true
		}
	}
	return false
}

func RemoteIP(conn *net.TCPConn) net.IP {
	return conn.RemoteAddr().(*net.TCPAddr).IP
}

func MustFindFreePort() int {
	p, err := freeport.GetFreePort()
	if err != nil {
		panic(err)
	}
	return p
}
