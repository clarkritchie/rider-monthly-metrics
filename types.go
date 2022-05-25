package main

import (
	"encoding/json"
	"fmt"
	"math"
)

const (
	METERS_TO_MILES = 0.000621371192
	METERS_TO_FEET  = 3.28084
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

func (d *Sum) ConvertSumToImperial() {
	d.Distance = d.Distance * METERS_TO_MILES
	d.Ascent = d.Ascent * METERS_TO_FEET
}

func (d *Max) ConvertMaxToImperial() {
	d.Distance = d.Distance * METERS_TO_MILES
	d.Ascent = d.Ascent * METERS_TO_FEET
}

func (d *Max) FormatMaxForEmail(m ...Max) map[string]interface{} {
	previousMax := Max{}
	// with veradic functions, the optional parameter is a slice
	if len(m) > 0 {
		previousMax = m[0]
	}
	data := make(map[string]interface{})

	// keep as float, but rounded
	data["distance"] = math.Round(d.Distance*100) / 100
	data["ascent"] = math.Round(d.Ascent*100) / 100

	time := make(map[string]interface{})
	hours, minutes := convertSecondsToHoursMinutes(d.Training) // d.Training is in seconds
	time["hours"] = hours
	time["minutes"] = minutes
	// compute the delta to the previous month
	if previousMax.Training > 0 {
		previousHours, previousMinutes := convertSecondsToHoursMinutes(previousMax.Training)
		time["hoursDelta"] = hours - previousHours
		time["minutesDelta"] = minutes - previousMinutes
	}
	data["training"] = time
	if previousMax.Distance > 0 {
		data["distanceDelta"] = d.Distance - previousMax.Distance
	}
	if previousMax.Ascent > 0 {
		data["ascentDelta"] = d.Ascent - previousMax.Ascent
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

	// keep as float, but rounded
	data["distance"] = math.Round(d.Distance*100) / 100
	data["ascent"] = math.Round(d.Ascent*100) / 100

	time := make(map[string]interface{})
	hours, minutes := convertSecondsToHoursMinutes(d.Training) // d.Training is in seconds
	time["hours"] = hours
	time["minutes"] = minutes
	// compute the delta to the previous month
	if previousSum.Training > 0 {
		previousHours, previousMinutes := convertSecondsToHoursMinutes(previousSum.Training)
		time["hoursDelta"] = hours - previousHours
		time["minutesDelta"] = minutes - previousMinutes
	}
	data["training"] = time
	if previousSum.Distance > 0 {
		data["distanceDelta"] = d.Distance - previousSum.Distance
	}
	if previousSum.Ascent > 0 {
		data["ascentDelta"] = d.Ascent - previousSum.Ascent
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

	data["current"] = current
	data["previous"] = previous

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
	fmt.Printf("monthlyData:\n")
	fmt.Printf("%+v", monthlyData)

	if monthlyData.Units == "imperial" {
		monthlyData.Previous.Sum.ConvertSumToImperial()
		monthlyData.Previous.Max.ConvertMaxToImperial()
		monthlyData.Current.Sum.ConvertSumToImperial()
		monthlyData.Current.Max.ConvertMaxToImperial()
	}

	emailData := monthlyData.FormatForEmail()
	return emailData
}
