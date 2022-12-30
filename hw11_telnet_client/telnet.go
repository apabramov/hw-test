package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Telnet struct {
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
	network string
	address string
	timeout time.Duration
}

func (t *Telnet) Connect() error {
	var err error
	if t.conn, err = net.DialTimeout(t.network, t.address, t.timeout); err != nil {
		return err
	}
	return nil
}

func (t *Telnet) Send() error {
	_, err := io.Copy(t.conn, t.in)
	return err
}

func (t Telnet) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	return err
}

func (t Telnet) Close() error {
	return t.conn.Close()
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Telnet{network: "tcp", address: address, timeout: timeout, in: in, out: out}
}
