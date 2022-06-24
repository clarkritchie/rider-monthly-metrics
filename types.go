package main

import (
	"encoding/json"
	"fmt"
	"math"
)

const (
	METERS_TO_MILES      = 0.000621371192
	METERS_TO_FEET       = 3.28084
	METERS_TO_KILOMETERS = 1000
)

type Sum struct {
	Training int     `json:"training"` // movement time, in seconds
	Distance float64 `json:"distance"` // meters
	Ascent   float64 `json:"ascent"`   // meters
}

type Max struct {
	Training int     `json:"training"` // movement time, in seconds
	Distance float64 `json:"distance"` // meters
	Ascent   float64 `json:"ascent"`   // meters
}

type MetricData struct {
	Sum   Sum    `json:"sum"`
	Max   Max    `json:"max"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type MonthlyData struct {
	Current  MetricData `json:"current"`
	Previous MetricData `json:"previous"`
	Units    string     `json:"units"`
}

// there is no built in Abs for integers
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func (d *Sum) ConvertSumToImperial() {
	d.Distance = d.Distance * METERS_TO_MILES
	d.Ascent = d.Ascent * METERS_TO_FEET
}

func (d *Max) ConvertMaxToImperial() {
	d.Distance = d.Distance * METERS_TO_MILES
	d.Ascent = d.Ascent * METERS_TO_FEET
}

func (d *Sum) ConvertSumDistanceToKms() {
	d.Distance = d.Distance / METERS_TO_KILOMETERS
}

func (d *Max) ConvertMaxDistanceToKms() {
	d.Distance = d.Distance / METERS_TO_KILOMETERS
}

func (d *Max) FormatMaxForEmail(m ...Max) map[string]interface{} {
	previousMax := Max{}
	// with veradic functions, the optional parameter is a slice
	if len(m) > 0 {
		previousMax = m[0]
	}
	data := make(map[string]interface{})

	data["distance"] = roundFloat(d.Distance, 1)
	data["ascent"] = roundFloat(d.Ascent, 1)

	time := make(map[string]interface{})
	hours, minutes := convertSecondsToHoursMinutes(d.Training) // d.Training is in seconds
	time["hours"] = hours
	time["minutes"] = minutes
	// compute the delta to the previous month
	if previousMax.Training >= 0 {
		previousHours, previousMinutes := convertSecondsToHoursMinutes(previousMax.Training)
		time["hoursDelta"] = abs(hours - previousHours)
		time["minutesDelta"] = abs(minutes - previousMinutes)
	}
	data["training"] = time
	if previousMax.Distance >= 0 {
		data["distanceDelta"] = roundFloat(math.Abs(d.Distance-previousMax.Distance), 1)
	}
	if previousMax.Ascent >= 0 {
		data["ascentDelta"] = roundFloat(math.Abs(d.Ascent-previousMax.Ascent), 1)
	}

	return data
}

func (d *Sum) FormatSumForEmail(s ...Sum) map[string]interface{} {
	previousSum := Sum{}
	// with veradic functions, the optional parameter is a slice
	if len(s) > 0 {
		previousSum = s[0]
	}
	data := make(map[string]interface{})

	data["distance"] = roundFloat(d.Distance, 1)
	data["ascent"] = roundFloat(d.Ascent, 1)

	time := make(map[string]interface{})
	hours, minutes := convertSecondsToHoursMinutes(d.Training) // d.Training is in seconds
	time["hours"] = hours
	time["minutes"] = minutes
	// compute the delta to the previous month
	if previousSum.Training >= 0 {
		previousHours, previousMinutes := convertSecondsToHoursMinutes(previousSum.Training)
		time["hoursDelta"] = abs(hours - previousHours)
		time["minutesDelta"] = abs(minutes - previousMinutes)
	}
	data["training"] = time
	if previousSum.Distance >= 0 {
		data["distanceDelta"] = roundFloat(math.Abs(d.Distance-previousSum.Distance), 1)
	}
	if previousSum.Ascent >= 0 {
		data["ascentDelta"] = roundFloat(math.Abs(d.Ascent-previousSum.Ascent), 1)
	}

	return data
}

func (d *MonthlyData) FormatForEmail() map[string]interface{} {
	data := make(map[string]interface{})
	current := make(map[string]interface{})
	previous := make(map[string]interface{})

	previous["sum"] = d.Previous.Sum.FormatSumForEmail()
	previous["max"] = d.Previous.Max.FormatMaxForEmail()
	previous["name"] = d.Previous.Name
	previous["count"] = d.Previous.Count

	current["sum"] = d.Current.Sum.FormatSumForEmail(d.Previous.Sum)
	current["max"] = d.Current.Max.FormatMaxForEmail(d.Previous.Max)
	current["name"] = d.Current.Name
	current["count"] = d.Current.Count
	current["countDelta"] = abs(d.Current.Count - d.Previous.Count)

	data["current"] = current
	data["previous"] = previous
	data["units"] = d.Units

	return data
}

func convertSecondsToHoursMinutes(seconds int) (int, int) {
	// 7200 seconds is 2 hours
	// 7260 seconds is 2 hours and 1 minute
	// 7270 seconds is still 2 hours and 1 minute -- the last 10 seconds are lost
	// 7319 seconds is also still 2 hours and 1 minute
	// 7320 seconds is 2 hours and 2 minutes
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	return hours, minutes
}

func ProcessMonthlyMetrics(data map[string]interface{}) map[string]interface{} {
	var monthlyData MonthlyData

	// read data back into a string (as byte array)
	dataAsString, _ := json.Marshal(data)
	fmt.Printf("ProcessMonthlyMetrics data: %+v", data)
	// now shove that back into the struct and let the hints do the field mapping magic
	if err := json.Unmarshal(dataAsString, &monthlyData); err != nil {
		fmt.Println(err)
	}
	if monthlyData.Units == "imperial" {
		monthlyData.Previous.Sum.ConvertSumToImperial()
		monthlyData.Previous.Max.ConvertMaxToImperial()
		monthlyData.Current.Sum.ConvertSumToImperial()
		monthlyData.Current.Max.ConvertMaxToImperial()
	} else {
		fmt.Println("conveting meters to kilometers")
		// these values are in meters, so make them kilometers
		monthlyData.Previous.Sum.ConvertSumDistanceToKms()
		monthlyData.Previous.Max.ConvertMaxDistanceToKms()
		monthlyData.Current.Sum.ConvertSumDistanceToKms()
		monthlyData.Current.Max.ConvertMaxDistanceToKms()
	}
	fmt.Printf("monthlyData:\n")
	fmt.Printf("%+v", monthlyData)
	emailData := monthlyData.FormatForEmail()
	return emailData
}
