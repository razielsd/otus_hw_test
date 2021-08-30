package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &HwTelnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type HwTelnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (h *HwTelnetClient) Connect() error {
	d := net.Dialer{
		Timeout: timeout,
	}
	con, err := d.Dial("tcp", h.address)
	if err != nil {
		return err
	}
	h.conn = con
	return nil
}

func (h *HwTelnetClient) Close() error {
	return h.conn.Close()
}

func (h *HwTelnetClient) Send() error {
	if _, err := io.Copy(h.conn, h.in); err != nil {
		return fmt.Errorf("error while sending: %w", err)
	}

	fmt.Fprintln(os.Stderr, "...EOF")

	return nil
}

func (h *HwTelnetClient) Receive() error {
	if _, err := io.Copy(h.out, h.conn); err != nil {
		return fmt.Errorf("error occurred while receiving: %w", err)
	}
	fmt.Fprintln(os.Stderr, "...Connection closed by remote host")
	return nil
}
