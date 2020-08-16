package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var host, port, timeout string

func init() {
	flag.StringVar(&timeout, "timeout", "10s", "timeout")
}

func worker(src io.Reader, fn func(string)) {
	scanner := bufio.NewScanner(src)
	for {
		if !scanner.Scan() {
			break
		}
		fn(scanner.Text())
	}
}

func main() {
	flag.Parse()
	tail := flag.Args()
	if len(tail) != 2 {
		log.Fatal("Please, provide host and port")
	}
	host, port = tail[0], tail[1]
	duration, err := time.ParseDuration(timeout)
	if err != nil {
		log.Fatal("Invalid duration")
	}

	in := &bytes.Buffer{}
	out := &bytes.Buffer{}

	tc := NewTelnetClient(net.JoinHostPort(host, port), duration, ioutil.NopCloser(in), out)
	if err = tc.Connect(); err != nil {
		os.Exit(1)
	}
	chDone := make(chan error)
	go worker(os.Stdin, func(text string) {
		in.WriteString(text + "\n")
		err := tc.Send()
		if err != nil {
			chDone <- err
		}
	})
	go worker(tc.GetConn(), func(text string) {
		_, err := io.Copy(os.Stdout, strings.NewReader(text+"\n"))
		if err != nil {
			chDone <- err
		}
	})

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-chDone:
			fmt.Fprintf(os.Stderr, "...%s\n", err)
			tc.Close()
			os.Exit(1)
		case <-termChan:
			fmt.Fprint(os.Stderr, "...Closing the connection")
			tc.Close()
			os.Exit(1)
		}
	}
}
