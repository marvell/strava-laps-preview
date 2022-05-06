package strava

import "time"

type ActivityType string

var (
	ActivityTypeRide ActivityType = "Ride"
	ActivityTypeRun  ActivityType = "Run"
)

type Activity struct {
	Id        int          `json:"id"`
	Name      string       `json:"name"`
	Type      ActivityType `json:"type"`
	StartDate time.Time    `json:"start_date"`
	Distance  float64      `json:"distance"`
}

type Lap struct {
	Id                 int       `json:"id"`
	Name               string    `json:"name"`
	ElapsedTime        int       `json:"elapsed_time"`
	MovingTime         int       `json:"moving_time"`
	StartDate          time.Time `json:"start_date"`
	StartDateLocal     time.Time `json:"start_date_local"`
	Distance           float64   `json:"distance"`
	StartIndex         int       `json:"start_index"`
	EndIndex           int       `json:"end_index"`
	TotalElevationGain float64   `json:"total_elevation_gain"`
	AverageSpeed       float64   `json:"average_speed"`
	MaxSpeed           float64   `json:"max_speed"`
	AverageCadence     float64   `json:"average_cadence"`
	// DeviceWatts        float64   `json:"device_watts"`
	// AverageWatts       float64   `json:"average_watts"`
	AverageHeartrate float64 `json:"average_heartrate"`
	MaxHeartrate     float64 `json:"max_heartrate"`
	LapIndex         int     `json:"lap_index"`
	Split            int     `json:"split"`
	PaceZone         int     `json:"pace_zone"`
}
