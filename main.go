package main

import (
	"encoding/json"
	"fmt"
	pp "github.com/k0kubun/pp/v3"
	"io/ioutil"
	"os"
)

//func roundFloat(val float64, precision uint) float64 {
//	ratio := math.Pow(10, float64(precision))
//	return math.Round(val*ratio) / ratio
//}

func processStuff(data map[string]interface{}, isMetric bool) {
	var monthlyData MonthlyData

	// read data back into a string (as byte array)
	detailAsString, _ := json.Marshal(data["details"])
	// now shove that back into the struct and let the hints do the field mapping magic
	if err := json.Unmarshal(detailAsString, &monthlyData); err != nil {
		fmt.Println(err)
	}

	//if !isMetric {
	//	monthlyData.Previous.Sum.ConvertSumToImperial()
	//	monthlyData.Previous.Max.ConvertMaxToImperial()
	//	monthlyData.Current.Sum.ConvertSumToImperial()
	//	monthlyData.Current.Max.ConvertMaxToImperial()
	//}

	emailData := monthlyData.FormatForEmail()
	pp.Print(emailData)
}

func main() {
	jsonFile, err := os.Open("data.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully opened data.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	var isMetric = true
	if result["units"] != "metric" {
		isMetric = false
	}
	processStuff(result, isMetric)
}

/*
func Oldmain() {
	jsonFile, err := os.Open("data.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened data.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	fmt.Println(result["details"])

	//// this essentially creates an instance of MonthlyData as data
	//data := &MonthlyData{
	//	Previous: MetricData{
	//		Sum: Sum{
	//			//Training: 86400, // 1 day, in seconds
	//			//Training: 3660, // 1:01, in seconds
	//			//Training: 3600,    // 1:00 seconds
	//			Training: 60, // 1:00 seconds
	//			// Distance: 20000.0, // meters
	//			Ascent: 100.0, // meters
	//		},
	//		Max: Max{
	//			Training: 1000,
	//			Distance: 1000.0,
	//			Ascent:   100.0,
	//		},
	//		Name:  "2020-08",
	//		Count: 2,
	//	},
	//	Current: MetricData{
	//		Sum: Sum{
	//			Training: 7200,
	//			Distance: 20000,
	//			Ascent:   100.0,
	//		},
	//		Max: Max{
	//			Training: 7200,
	//			Distance: 20000,
	//			Ascent:   100.0,
	//		},
	//		Name:  "2020-09",
	//		Count: 1,
	//	},
	//}
	//
	//fmt.Println("raw metric data:")
	//pp.Print(data)
	//
	//// convert units
	//// this is the equivalent of calling a method on an instance of a class
	//// however, instead of class data is a struct, Previous is also a struct, and Sum is also a struct
	//// the Sum struct has a receiver function called ConvertSumToImperial that lets me call it
	//// like a method on an instance of a class
	//data.Previous.Sum.ConvertSumToImperial()
	//data.Previous.Max.ConvertMaxToImperial()
	//data.Current.Sum.ConvertSumToImperial()
	//data.Current.Max.ConvertMaxToImperial()
	//
	//fmt.Println("\nconverted to imperial:")
	//pp.Print(data)
	//
	//fmt.Printf("\nFinal, formatted values:\n")
	//pp.Printf("Sum Training goes from %s to %s\n", data.Previous.Sum.Training, convertSecondsToTimeString(data.Previous.Sum.Training))
	//
	//roundedDistance := roundFloat(data.Previous.Sum.Distance, 1)
	//distanceAsString := fmt.Sprintf("%v", roundedDistance)
	//pp.Printf("Sum Distance goes from %s to %s\n", data.Previous.Sum.Distance, distanceAsString)
	//
	//// round then format ascent as string
	//roundedAscent := roundFloat(data.Previous.Sum.Ascent, 0)
	//ascentAsString := fmt.Sprintf("%v", roundedAscent)
	//pp.Printf("Sum Ascent goes from %s to %v\n", data.Previous.Sum.Ascent, ascentAsString)
}
*/
