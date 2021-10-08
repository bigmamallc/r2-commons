package piper

import (
	"bytes"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/bigmamallc/r2-commons/api/tcp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"io"
	"math/rand"
	"net"
	"os"
	"sync"
	"testing"
	"time"
)

func startEchoServer(port int, log zerolog.Logger) {
	addr := &net.TCPAddr{Port: port}
	lsnr, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		panic(err)
	}
	go func() {
		defer lsnr.Close()

		for {
			if err := lsnr.SetDeadline(time.Now().Add(time.Second)); err != nil {
				panic(err)
			}
			conn, err := lsnr.AcceptTCP()
			if os.IsTimeout(err) {
				continue
			} else if err != nil {
				panic(err)
			}
			go func() {
				defer conn.Close()

				buf := make([]byte, 256)
				for {
					readSize, readErr := conn.Read(buf)
					log.Debug().Err(readErr).Int("readSize", readSize).Msg("conn.Read()")
					var writeErr error
					if readSize > 0 {
						var writeSize int
						writeSize, writeErr = conn.Write(buf[:readSize])
						log.Debug().Err(writeErr).Int("writeSize", writeSize).Msg("conn.Write()")
					}
					if readErr != nil || writeErr != nil {
						return
					}
				}
			}()
		}
	}()
	if !tcp.Await(addr, time.Now().Add(time.Second)) {
		panic("server did not start")
	}
}

func TestPiper(t *testing.T) {
	log := zerolog.New(zerolog.NewConsoleWriter()).Level(zerolog.InfoLevel)
	mxReg := prometheus.NewRegistry()

	p := New(log, mxReg, "test", 1024, time.Second)

	port1 := tcp.MustFindFreePort()
	port2 := tcp.MustFindFreePort()

	startEchoServer(port1, log.With().Str("server", "echo1").Logger())
	startEchoServer(port2, log.With().Str("server", "echo2").Logger())

	conn1, err := net.DialTCP("tcp4", nil, &net.TCPAddr{Port: port1})
	if err != nil { panic(err) }
	defer func() {
		log.Info().Msg("closing conn1")
		conn1.Close()
	}()

	conn2, err := net.DialTCP("tcp4", nil, &net.TCPAddr{Port: port2})
	if err != nil { panic(err) }
	defer func() {
		log.Info().Msg("closing conn2")
		conn2.Close()
	}()

	p.StartBidiPipe(conn1, conn2)

	sendAndReceive := func(a *net.TCPConn, b *net.TCPConn, maxMsgSize int) {
		msg := []byte(randomdata.Alphanumeric(rand.Intn(maxMsgSize-1)+1))
		if _, err := a.Write(msg); err != nil {
			panic(err)
		}
		recvMsg := make([]byte, len(msg))
		if _, err := io.ReadFull(b, msg); err != nil {
			if !bytes.Equal(msg, recvMsg) {
				panic(fmt.Sprintf("msg: %s recvMsg: %s", string(msg), string(recvMsg)))
			}
		}
	}

	const N = 100
	const maxDelayMS = 100
	const maxMsgSize = 1024

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			time.Sleep(time.Millisecond * time.Duration(rand.Int() % maxDelayMS))
			log.Debug().Str("dir", "1->2").Msg("sendAndReceive()")
			sendAndReceive(conn1, conn2, maxMsgSize)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < N; i ++ {
			time.Sleep(time.Millisecond * time.Duration(rand.Int() % maxDelayMS))
			log.Debug().Str("dir", "2->1").Msg("sendAndReceive()")
			sendAndReceive(conn2, conn1, maxMsgSize)
		}
	}()
	wg.Wait()
}
