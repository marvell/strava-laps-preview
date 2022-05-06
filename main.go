package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"github.com/marvell/strava-laps-preview/strava"
)

var (
	interval = 5 * time.Minute
	limit    = 10
)

var (
	flagSocks5Addr string
	flagSocks5User string
	flagSocks5Pass string

	flagStravaApiClientId         string
	flagStravaApiClientSecret     string
	flagStravaApiRefreshToken     string
	flagStravaStartFromActivityId int

	flagDebugMode bool
)

func init() {
	flag.StringVar(&flagSocks5Addr, "socks5-addr", "", "Socks5 address")
	flag.StringVar(&flagSocks5User, "socks5-user", "", "Socks5 user")
	flag.StringVar(&flagSocks5Pass, "socks5-pass", "", "Socks5 password")

	flag.StringVar(&flagStravaApiClientId, "client-id", "", "Strava API client ID")
	flag.StringVar(&flagStravaApiClientSecret, "client-secret", "", "Strava API client secret")
	flag.StringVar(&flagStravaApiRefreshToken, "refresh-token", "", "Strava API refresh token")
	flag.IntVar(&flagStravaStartFromActivityId, "start-from-id", 0, "Strava activity ID start from")

	flag.BoolVar(&flagDebugMode, "debug", false, "Debug mode")
}

func main() {
	flag.Parse()

	var opts []strava.Option
	if flagSocks5Addr != "" {
		opts = append(opts, strava.WithSocks5(flagSocks5Addr, flagSocks5User, flagSocks5Pass))
	}
	if flagDebugMode {
		opts = append(opts, strava.WithDebugMode())
	}

	c, err := strava.NewClient(flagStravaApiClientId, flagStravaApiClientSecret, flagStravaApiRefreshToken, opts...)
	if err != nil {
		log.Fatal(err)
	}

	var lastActivityId = flagStravaStartFromActivityId

	for range FirstTicker(interval) {
		log.Print("get activity list")
		activities, err := c.GetAthleteActivities(limit)
		if err != nil {
			log.Fatalf("c.GetAthleteActivities: %#v", err)
		}

		if len(activities) == 0 {
			log.Print("there are no activities")
			continue
		}

		sort.Slice(activities, func(i, j int) bool { return activities[i].StartDate.Before(activities[j].StartDate) })

		if lastActivityId == 0 {
			lastActivityId = activities[len(activities)-1].Id
			log.Printf("mark %d as last activity, nothing to update", lastActivityId)
			continue
		}

		for _, a := range activities {
			if a.Type != strava.ActivityTypeRun {
				log.Printf("%d has unsuitable type: %q, skip", a.Id, a.Type)
				continue
			}

			if a.Id <= lastActivityId {
				continue
			}

			log.Printf("[%d] try to update", a.Id)
			updateActivity(a.Id, c)

			lastActivityId = a.Id
		}

		log.Printf("wait for %s", interval)
	}
}

func updateActivity(id int, c *strava.Client) {
	laps, err := c.GetActivityLaps(id)
	if err != nil {
		log.Fatalf("c.GetActivityLaps: %#v", err)
	}

	if len(laps) < 2 {
		log.Printf("[%d] has only a lap, nothing to update", id)
		return
	}

	var desc string
	for i, l := range laps {
		avgPace := ConvertSpeedToPace(l.AverageSpeed)
		avgHr := int(math.Round(l.AverageHeartrate))

		desc += fmt.Sprintf("%s %d) %s / %s / %s / %d\n", speedEmoji(l.AverageSpeed), i+1, FormatDistance(l.Distance),
			FormatDuration(l.MovingTime), FormatDuration(avgPace), avgHr)
	}

	err = c.UpdateActivityDescription(id, desc)
	if err != nil {
		log.Fatalf("c.UpdateActivityDescription: %#v", err)
	}

	log.Printf("[%d] has updated", id)
}

var speedToEmoji = []struct {
	maxSpeed float64
	emoji    string
}{
	{10.344828, "🟣"}, // 5:48
	{11.076923, "🔵"}, // 5:25
	{11.881188, "🟢"}, // 5:03 - Threshold
	{12.811388, "🟡"}, // 4:41
	{13.953488, "🟠"}, // 4:18
	{30.0, "🔴"},
}

func speedEmoji(speed float64) string {
	speed = speed * 3.6

	for _, s := range speedToEmoji {
		if speed <= s.maxSpeed {
			return s.emoji
		}
	}
	return speedToEmoji[len(speedToEmoji)-1].emoji
}
