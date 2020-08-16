package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var (
	ErrConnectionClosed  = errors.New("...Connection closed by peer")
	ErrConnectionRefused = errors.New("...Connection refused")
)

type TelnetClient interface {
	Connect() error
	GetConn() net.Conn
	Close() error
	Send() error
	Receive() error
}

type Client struct {
	address string
	conn    net.Conn
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func (tc *Client) GetConn() net.Conn {
	return tc.conn
}

func (tc *Client) Connect() error {
	dialer := &net.Dialer{}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, tc.timeout)
	defer cancel()

	conn, err := dialer.DialContext(ctx, "tcp", tc.address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "...Cannot connect: %v\n", ErrConnectionRefused)

		return err
	}
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", tc.address)

	tc.conn = conn

	return nil
}

func (tc *Client) Close() error {
	if tc.conn == nil {
		return nil
	}
	err := tc.conn.Close()
	if err != nil {
		return ErrConnectionClosed
	}

	return nil
}

func (tc *Client) Send() error {
	_, err := bufio.NewWriter(tc.conn).ReadFrom(tc.in)
	if err != nil {
		return ErrConnectionClosed
	}

	return err
}

func (tc *Client) Receive() error {
	_, err := bufio.NewWriter(tc.out).ReadFrom(tc.conn)

	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{address: address, timeout: timeout, in: in, out: out}
}
