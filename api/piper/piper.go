package piper

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"io"
	"net"
	"time"
)

type Piper struct {
	log zerolog.Logger

	bufferSize   int
	writeTimeout time.Duration

	activeUniPipes prometheus.Gauge
	totalBytesAtoB prometheus.Counter
	totalBytesBtoA prometheus.Counter
	ioErrors       prometheus.Counter
}

func New(log zerolog.Logger, mxReg *prometheus.Registry, mxSubsystem string,
	bufferSize int, writeTimeout time.Duration) *Piper {
	p := &Piper{
		log: log,

		bufferSize:   bufferSize,
		writeTimeout: writeTimeout,

		activeUniPipes: prometheus.NewGauge(prometheus.GaugeOpts{Name: "activeUniPipes", Subsystem: mxSubsystem}),
		totalBytesAtoB: prometheus.NewCounter(prometheus.CounterOpts{Name: "totalBytesAtoB", Subsystem: mxSubsystem}),
		totalBytesBtoA: prometheus.NewCounter(prometheus.CounterOpts{Name: "totalBytesBtoA", Subsystem: mxSubsystem}),
		ioErrors:       prometheus.NewCounter(prometheus.CounterOpts{Name: "ioErrors", Subsystem: mxSubsystem}),
	}
	mxReg.MustRegister(p.activeUniPipes, p.totalBytesAtoB, p.totalBytesBtoA, p.ioErrors)
	return p
}

func (p *Piper) StartBidiPipe(connA net.Conn, connB net.Conn) {
	p.StartBidiPipeWithLog(connA, connB, p.log)
}

func (p *Piper) StartBidiPipeWithLog(connA net.Conn, connB net.Conn, baseLog zerolog.Logger) {
	abLog := baseLog
	baLog := baseLog
	if p.log.GetLevel() <= zerolog.DebugLevel {
		abLog = abLog.With().Str("dir", "a->b").Logger()
		baLog = baLog.With().Str("dir", "b->a").Logger()
	}

	go p.uniPipe(connA, connB, abLog, p.totalBytesAtoB)
	go p.uniPipe(connB, connA, baLog, p.totalBytesBtoA)
}

func (p *Piper) uniPipe(from net.Conn, to net.Conn, log zerolog.Logger, byteMetric prometheus.Counter) {
	log.Debug().Msg("entering uniPipe()")

	p.activeUniPipes.Inc()

	defer func() {
		p.activeUniPipes.Dec()

		from.Close()
		to.Close()
		log.Debug().Msg("exiting uniPipe()")
	}()

	if err := from.SetReadDeadline(time.Time{}); err != nil {
		log.Warn().Err(err).Msg("from.SetReadDeadline() failed")
		return
	}
	buf := make([]byte, p.bufferSize)
	for {
		readSize, readErr := from.Read(buf)
		log.Info().Err(readErr).Int("readSize", readSize).Msg("from.Read()")
		if readSize > 0 {
			if err := to.SetWriteDeadline(time.Now().Add(p.writeTimeout)); err != nil {
				log.Warn().Err(err).Msg("to.SetWriteDeadline() failed")
				return
			}
			writeSize, writeErr := to.Write(buf[:readSize])
			log.Info().Err(writeErr).Int("writeSize", writeSize).Msg("to.Write()")
			if writeErr != nil {
				log.Warn().Err(writeErr).Msg("to.Write() failed")
				return
			}

			byteMetric.Add(float64(writeSize))
		}
		if errors.Is(readErr, io.EOF) {
			log.Debug().Msg("io.EOF from from.Read() - closing from and exiting")
			return
		} else if readErr != nil {
			log.Debug().Err(readErr).Msg("error from from.Read() - exiting")
			return
		}
		log.Debug().Msg("proceeding to the next read/write cycle")
	}
}
