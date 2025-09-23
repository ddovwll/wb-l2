package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"qwe/internal/sky_scrapper"
	"runtime"
	"syscall"
	"time"
)

func main() {
	var depth uint
	var folder string
	var workers uint
	var timeout uint
	var httpTimeout uint

	flag.UintVar(&depth, "depth", 1, "max recursion depth")
	flag.StringVar(&folder, "folder", "content", "content folder")
	flag.UintVar(&workers, "workers", uint(runtime.GOMAXPROCS(0)), "workers limit")
	flag.UintVar(&timeout, "timeout", 1, "app timeout in minutes")
	flag.UintVar(&httpTimeout, "httpTimeout", 10, "http timeout in seconds")

	flag.Parse()

	args := flag.Args()

	url := args[0]

	if url == "" {
		log.Fatal("url is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*time.Duration(timeout))
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		cancel()
		<-sig
		os.Exit(1)
	}()

	scrapper, err := sky_scrapper.NewScrapper(ctx, folder, depth, int(httpTimeout))
	if err != nil {
		log.Fatal(err)
	}

	err = scrapper.Start(url, int(workers))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("done")
}
