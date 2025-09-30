package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"task17/internal/telnet"
	"time"
)

func main() {
	var (
		host    string
		port    int
		timeout int
	)

	flag.IntVar(&timeout, "timeout", 10, "connection timeout")
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		log.Fatal("incorrect arguments")
	}

	host = args[0]
	port, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	err = telnet.Run(ctx, host, port, time.Duration(timeout)*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()
}
