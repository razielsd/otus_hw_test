package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
}

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal("require arguments: host, port")
	}
	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatal("Unable connect", err)
	}
	defer client.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	go func() {
		defer cancel()
		err := client.Send()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	go func() {
		defer cancel()
		err := client.Receive()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	<-ctx.Done()
}
