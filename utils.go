package main

import (
	"fmt"
	"math"
	"time"
)

func FormatDistance(d float64) string {
	d = math.Round(d/10) * 10

	if d >= 1000 {
		m := d / 1000

		if int(d)%1000 == 0 {
			return fmt.Sprintf("%.0fkm", m)
		}

		if int(d)%100 == 0 {
			return fmt.Sprintf("%.1fkm", m)
		}

		return fmt.Sprintf("%.2fkm", m)
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

func PaceToSpeed(pace time.Duration) float64 {
	return 3600. / pace.Seconds()
}
