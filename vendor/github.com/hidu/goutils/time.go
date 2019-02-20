package utils
import "github.com/equalll/mydebug"

import (
	"time"
)

func SetInterval(call func(), sec int64) *time.Ticker {mydebug.INFO()
	ticker := time.NewTicker(time.Duration(sec) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				call()
			}
		}
	}()
	return ticker
}
