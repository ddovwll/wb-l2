package time_scrapper

import (
	"github.com/beevik/ntp"
	"time"
)

const timeUrl = "0.beevik-ntp.pool.ntp.org"

func GetNetworkTime() (time.Time, error) {
	time, err := ntp.Time(timeUrl)
	return time, err
}

func GetNetworkTimeByUrl(url string) (time.Time, error) {
	time, err := ntp.Time(url)
	return time, err
}
