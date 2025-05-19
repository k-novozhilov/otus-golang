package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "timeout for connection")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go-telnet [--timeout=10s] host port")
		os.Exit(1)
	}

	address := net.JoinHostPort(args[0], args[1])
	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	done := make(chan struct{})

	go func() {
		<-ctx.Done()
		fmt.Fprintln(os.Stderr, "...Connection terminated")
		client.Close()
		os.Exit(0)
	}()

	go func() {
		err := client.Receive()
		if err != nil {
			fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
			cancel()
		}
		close(done)
	}()

	err = client.Send()
	if err != nil {
		fmt.Fprintln(os.Stderr, "...EOF")
		cancel()
	}

	<-done
	client.Close()
}
