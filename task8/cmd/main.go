package main

import (
	"log"
	"task8/pkg/time_scrapper"
)

func main() {
	time, err := time_scrapper.GetNetworkTime()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Current time from NTP server:", time.String())
}
