package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", time.Second*10, "Timeout connection")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprint(os.Stderr, "error address.")
		os.Exit(1)
	}

	address := net.JoinHostPort(args[0], args[1])

	t := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := t.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "connect error: %v", err)
		return
	}
	defer t.Close()
	fmt.Fprintf(os.Stderr, "...Connected to %v", address)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		if err := t.Send(); err != nil {
			fmt.Fprintf(os.Stderr, "send error: %v", err)
		} else {
			fmt.Fprintf(os.Stderr, "...EOF")
		}
		cancel()
	}()

	go func() {
		if err := t.Receive(); err != nil {
			fmt.Fprintf(os.Stderr, "receive error: %v", err)
		} else {
			fmt.Fprint(os.Stderr, "...Connection was closed by peer")
		}
		cancel()
	}()

	<-ctx.Done()
}
