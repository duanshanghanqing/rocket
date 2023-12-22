package utils

import (
	"testing"
	"time"
)

// go test -v -run Test_NewTicker ticker_test.go ticker.go
func Test_NewTicker(t *testing.T) {
	ticker := NewTicker(time.Second, func() {
		t.Log("Ticker fired")
	})
	ticker.Start()

	time.Sleep(10 * time.Second)
	ticker.Stop()
	time.Sleep(time.Second)
	t.Log("Ticker Stop")
}
