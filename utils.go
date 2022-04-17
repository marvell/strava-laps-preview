package main

import (
	"fmt"
	"math"
	"time"
)

func FormatDistance(d float64) string {
	m := int(math.Round(d))

	if m >= 1000 {
		m = int(m/10) * 10 // round to 10

		if m%1000 == 0 {
			return fmt.Sprintf("%.0fkm", d/1000)
		}

		if m%100 != 0 {
			return fmt.Sprintf("%.2fkm", d/1000)
		}

		return fmt.Sprintf("%.1fkm", d/1000)
	}

	return fmt.Sprintf("%.0fm", d)
}

func FormatDuration(secs int) string {
	h := secs / 3600
	m := secs % 3600 / 60
	s := secs % 60

	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}

	return fmt.Sprintf("%02d:%02d", m, s)
}

func ConvertSpeedToPace(s float64) int {
	return int(math.Round(1000 / s))
}

func FirstTicker(d time.Duration) <-chan time.Time {
	ch := make(chan time.Time, 1)
	ch <- time.Now()

	ticker := time.NewTicker(d)

	go func() {
		for t := range ticker.C {
			ch <- t
		}
	}()

	return ch
}
