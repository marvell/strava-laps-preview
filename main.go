package main

import (
	"flag"
	"fmt"
	"log"
	"math"

	"github.com/marvell/strava-laps-preview/strava"
)

var (
	flagSocks5Addr       string
	flagSocks5User       string
	flagSocks5Pass       string
	flagStravaApiToken   string
	flagStravaActivityId int
	flagDebugMode        bool
)

func init() {
	flag.StringVar(&flagSocks5Addr, "socks5-addr", "", "Socks5 address")
	flag.StringVar(&flagSocks5User, "socks5-user", "", "Socks5 user")
	flag.StringVar(&flagSocks5Pass, "socks5-pass", "", "Socks5 password")
	flag.StringVar(&flagStravaApiToken, "token", "", "Strava API token")
	flag.IntVar(&flagStravaActivityId, "activity-id", 0, "Strava activity ID")
	flag.BoolVar(&flagDebugMode, "debug", false, "Debug mode")
}

func main() {
	flag.Parse()

	c, err := strava.NewClient(flagStravaApiToken,
		strava.WithSocks5(flagSocks5Addr, flagSocks5User, flagSocks5Pass),
		strava.WithDebugMode(flagDebugMode))
	if err != nil {
		log.Fatal(err)
	}

	laps, err := c.GetActivityLaps(flagStravaActivityId)
	if err != nil {
		log.Fatalf("c.GetActivityLaps: %#v", err)
	}

	if len(laps) < 2 {
		log.Print("only 1 lap, nothing to update")
		return
	}

	var desc string
	for i, l := range laps {
		avgPace := ConvertSpeedToPace(l.AverageSpeed)
		avgHr := int(math.Round(l.AverageHeartrate))

		desc += fmt.Sprintf("%d) %s / %s / %s / %d\n", i+1, FormatDistance(l.Distance), FormatDuration(l.MovingTime), FormatDuration(avgPace), avgHr)
	}

	err = c.UpdateActivityDescription(flagStravaActivityId, desc)
	if err != nil {
		log.Fatalf("c.UpdateActivityDescription: %#v", err)
	}

	log.Printf("%d activity has updated", flagStravaActivityId)
}
