package mqtt

import (
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/errors"
	"github.com/mainflux/mproxy/pkg/session"
)

const (
	network   = "udp"
	maxPktLen = 1500
)

// Proxy is main MQTT proxy struct
type Proxy struct {
	address string
	target  string
	handler session.Handler
	logger  logger.Logger
	dialer  net.Dialer
}

// New returns a new mqtt Proxy instance.
func New(address, target string, handler session.Handler, logger logger.Logger) *Proxy {
	return &Proxy{
		address: address,
		target:  target,
		handler: handler,
		logger:  logger,
	}
}

func (p Proxy) accept(l *net.UDPConn) {
	buff := make([]byte, maxPktLen)

	for {
		nr, _, err := l.ReadFromUDP(buff)
		if err != nil {
			if neterr, ok := err.(net.Error); ok && (neterr.Temporary() || neterr.Timeout()) {
				time.Sleep(5 * time.Millisecond)
				continue
			}
			// return err
		}
		tmp := make([]byte, nr)
		copy(tmp, buff)
		p.handle(l)
		// go handlePacket(listener, tmp, addr, rh)

		// p.logger.Info("Accepted new client")
		// go p.handle(conn)
	}
}

func (p Proxy) handle(inbound net.Conn) {
	defer p.close(inbound)
	uaddr, err := net.ResolveUDPAddr(network, p.target)
	if err != nil {
		p.logger.Warn("Unable to resolve target address" + p.target)
	}

	outbound, err := net.DialUDP(network, nil, uaddr)
	defer p.close(outbound)

	s := session.New(inbound, outbound, p.handler, p.logger, x509.Certificate{})

	if err = s.Stream(); !errors.Contains(err, io.EOF) {
		p.logger.Warn("Broken connection for client: " + s.Client.ID + " with error: " + err.Error())
	}
}

// Proxy of the server, this will block.
func (p Proxy) Proxy() error {
	uaddr, err := net.ResolveUDPAddr(network, p.address)
	if err != nil {
		return err
	}

	l, err := net.ListenUDP(network, uaddr)
	if err != nil {
		return err
	}

	// Acceptor loop
	p.accept(l)

	p.logger.Info("Server Exiting...")
	return nil
}

func (p Proxy) close(conn net.Conn) {
	if err := conn.Close(); err != nil {
		p.logger.Warn(fmt.Sprintf("Error closing connection %s", err.Error()))
	}
}
