package utils

import (
	"time"
)

type Ticker struct {
	ticker *time.Ticker
	fn     func()
}

func NewTicker(d time.Duration, fn func()) *Ticker {
	return &Ticker{
		ticker: time.NewTicker(d),
		fn:     fn,
	}
}

func (t *Ticker) Start() {
	go func() {
		for range t.ticker.C {
			t.fn()
		}
	}()
}

func (t *Ticker) Stop() {
	t.ticker.Stop()
}

//func main() {
//	ticker := NewTicker(1*time.Second, func() {
//		fmt.Println("Ticker fired")
//	})
//	ticker.Start()
//
//	time.Sleep(10 * time.Second)
//	ticker.Stop()
//}

//func main() {
//	ticker := time.NewTicker(1 * time.Second)
//
//	go func() {
//		for range ticker.C {
//			fmt.Println("Ticker fired")
//		}
//	}()
//
//	time.Sleep(10 * time.Second)
//	ticker.Stop()
//
//	fmt.Println("Ticker stopped")
//}
