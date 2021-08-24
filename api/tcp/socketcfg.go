package tcp

import (
	"fmt"
	"net"
	"time"
)

type SocketCfg struct {
	KeepAlive       bool          `env:"KEEP_ALIVE" default:"true"`
	KeepAlivePeriod time.Duration `env:"KEEP_ALIVE_PERIOD" default:"0"`
	NoDelay         bool          `env:"NO_DELAY" default:"true"`
	ReadBuffer      int           `env:"READ_BUFFER" default:"0"`
	WriteBuffer     int           `env:"WRITE_BUFFER" default:"0"`
}

func ApplySocketCfg(conn *net.TCPConn, cfg *SocketCfg) error {
	if err := conn.SetKeepAlive(cfg.KeepAlive); err != nil {
		return fmt.Errorf("SetKeepAlive() failed: %w", err)
	}
	if cfg.KeepAlivePeriod > 0 {
		if err := conn.SetKeepAlivePeriod(cfg.KeepAlivePeriod); err != nil {
			return fmt.Errorf("SetKeepAlivePeriod() failed: %w", err)
		}
	}
	if err := conn.SetNoDelay(cfg.NoDelay); err != nil {
		return fmt.Errorf("SetNoDelay() failed: %w", err)
	}
	if cfg.ReadBuffer > 0 {
		if err := conn.SetReadBuffer(cfg.ReadBuffer); err != nil {
			return fmt.Errorf("SetReadBuffer() failed: %w", err)
		}
	}
	if cfg.WriteBuffer > 0 {
		if err := conn.SetWriteBuffer(cfg.WriteBuffer); err != nil {
			return fmt.Errorf("SetWriteBuffer() failed: %w", err)
		}
	}
	return nil
}
