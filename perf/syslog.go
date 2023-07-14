package perf

import (
	"github.com/schollz/progressbar/v3"
	"golang.org/x/time/rate"
	"log/syslog"
	"math/rand"
	"sync"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func Run(max int, workers int, qps int, timeout int, host string, tag string, messageSize int) {
	var wg sync.WaitGroup

	duration := time.Duration(timeout) * time.Second
	bar := progressbar.Default(int64(max))

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			runWorker(max/workers, qps/workers, duration, bar, host, tag, messageSize)
			wg.Done()
		}()
	}
	wg.Wait()
}

func runWorker(max int, qps int, timeout time.Duration, bar *progressbar.ProgressBar, host string, tag string, messageSize int) {
	if max < 1 {
		return
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	LOGGER, _ := syslog.Dial("udp", host, syslog.LOG_INFO|syslog.LOG_AUTH, tag)
	message := RandStringBytes(messageSize)

	n := rate.Every(time.Second / time.Duration(qps))
	limiter := rate.NewLimiter(n, 1)

	done := make(chan bool, 1)
	i := 0

	for {
		select {
		case <-timer.C:
			return
		case <-done:
			return
		default:
			if limiter.Allow() {
				i++
				_ = LOGGER.Info(message)
				_ = bar.Add(1)
			}
			if i == max {
				done <- true
			}
		}
	}
}
