package main

import (
	"errors"
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

type Client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", c.address)
	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) Send() error {
	_, err := io.Copy(c.conn, c.in)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}

func (c *Client) Receive() error {
	_, err := io.Copy(c.out, c.conn)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}
