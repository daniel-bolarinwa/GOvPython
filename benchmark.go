package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"sync"

	//	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"

	"github.com/gonum/matrix/mat64"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/ensemble"
	"github.com/sjwhitworth/golearn/evaluation"
)

func main() {
	file, err := os.Open("./flights.csv")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	flights := dataframe.ReadCSV(file)

	file, err = os.Open("./airports.csv")
	if err != nil {
		log.Fatal(err)
	}
	airports := dataframe.ReadCSV(file)

	fmt.Println(flights)
	fmt.Println(airports)

	fmt.Println(flights.Names())

	flights = flights.Drop([]int{0, 5, 6, 10, 12, 13, 14, 15, 16, 18, 19, 20, 21, 23, 24, 25, 26, 27, 28, 29, 30})

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	origAirport := flights.Col("ORIGIN_AIRPORT").Records()
    for i, val := range origAirport {
		if _, err := strconv.ParseFloat(val, 64); err == nil {
			origAirport[i] = "OTHER"
		}
    }
    flights = flights.Mutate(series.New(origAirport, series.String, "ORIGIN_AIRPORT"))

	destAirport := flights.Col("DESTINATION_AIRPORT").Records()
    for i, val := range destAirport {
		if _, err := strconv.ParseFloat(val, 64); err == nil {
			destAirport[i] = "OTHER"
		}
    }
	flights = flights.Mutate(series.New(destAirport, series.String, "DESTINATION_AIRPORT"))

	records := flights.Records()
	
	rowsToRemove := make(map[int]bool)
	for idx, val := range records {
		for _, v := range val {
			if v == "NaN" {
				rowsToRemove[idx] = true
			}
		}
	}

	j := 0
    for i := 0; i < len(records); i++ {
        if !rowsToRemove[i] {
            records[j] = records[i]
            j++
        }
    }
    records = records[:j]

	flights = dataframe.LoadRecords(records)
	dayOfWeek := flights.Col("DAY_OF_WEEK").Records()
    for i, val := range dayOfWeek {
		switch val {
		case "7":
			dayOfWeek[i] = "SUNDAY"
		case "1":
			dayOfWeek[i] = "MONDAY"
		case "2":
			dayOfWeek[i] = "TUESDAY"
		case "3":
			dayOfWeek[i] = "WEDNESDAY"
		case "4":
			dayOfWeek[i] = "THURSDAY"
		case "5":
			dayOfWeek[i] = "FRIDAY"
		case "6":
			dayOfWeek[i] = "SATURDAY"
		}
    }
	flights = flights.Mutate(series.New(dayOfWeek, series.String, "DAY_OF_WEEK"))

	flights = OneHotCode("AIRLINE", flights)
	flights = OneHotCode("ORIGIN_AIRPORT", flights)
	flights = OneHotCode("DESTINATION_AIRPORT", flights)
	flights = OneHotCode("DAY_OF_WEEK", flights)

	fileOut, err := os.Create("output.csv")
	if err != nil {
		log.Fatal(err)
	}

	flights.WriteCSV(fileOut)

	data, err := base.ParseCSVToInstances("output.csv", true)
	if err != nil {
		log.Fatal(err)
	}

	rf := ensemble.NewRandomForest(100, 2)
	cv, err := evaluation.GenerateCrossFoldValidationConfusionMatrices(data, rf, 4)
	if err != nil {
		log.Fatal(err)
	}

	// Get the mean, variance and standard deviation of the accuracy for the
	// cross validation.
	mean, variance := evaluation.GetCrossValidatedMetric(cv, evaluation.GetAccuracy)
	stdev := math.Sqrt(variance)

	// Output the cross metrics to standard out.
	fmt.Printf("\nAccuracy\n%.2f (+/- %.2f)\n\n", mean, stdev*2)

	// start := time.Now()
	// elapsed := time.Since(start)

}

// One Hot Code
func OneHotCode(columnName string, dataframe dataframe.DataFrame) dataframe.DataFrame {
    col := dataframe.Col(columnName)
    values := make(map[string]bool)
    for i := 0; i < col.Len(); i++ {
        values[col.Elem(i).String()] = true
    }

    uniqueValues := make([]string, len(values))
    j := 0
    for value := range values {
        uniqueValues[j] = value
        j++
    }

	// Create a binary matrix for the one-hot encoding
    m := mat64.NewDense(len(dataframe.Records())-1, len(uniqueValues), nil)
	var wg sync.WaitGroup
	for j, uniqueValue := range uniqueValues {
		wg.Add(1)
		go func(j int, uniqueValue string) {
			defer wg.Done()
			col := dataframe.Col(columnName)
			for i := 0; i < col.Len(); i++ {
				value := col.Elem(i).String()
				if value == uniqueValue {
					m.Set(i, j, 1)
				}
			}
		}(j, uniqueValue)
	}
	wg.Wait()

	var wg1 sync.WaitGroup
	wg1.Add(len(uniqueValues))
	for j, uniqueValue := range uniqueValues {
		go func(j int, uniqueValue string) {
			colView := m.ColView(j)
			newColData := make([]int, len(dataframe.Records())-1)
			for i := 0; i < len(newColData); i++ {
				newColData[i] = int(colView.At(i, 0))
			}
    		dataframe = dataframe.Mutate(series.New(newColData, series.Int, columnName + "_" + uniqueValue))
			wg1.Done()
		}(j, uniqueValue)
	}
	wg1.Wait()
    
    dataframe = dataframe.Drop(columnName)

	return dataframe
}
