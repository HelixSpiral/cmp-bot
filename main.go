package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/HelixSpiral/cmpscraper"
)

func main() {
	stats, err := cmpscraper.GetStats()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\r\n", stats)

	statList := buildMessages(stats)

	err = sendMessages(statList)
	if err != nil {
		panic(err)
	}
}

func buildMessages(stats cmpscraper.CMP) []string {
	var statList []string

	totalCust, err := strconv.ParseFloat(strings.ReplaceAll(stats.Total, ",", ""), 64)
	if err != nil {
		panic(err)
	}

	totalOut, err := strconv.ParseFloat(strings.ReplaceAll(stats.WithoutPower, ",", ""), 64)
	if err != nil {
		panic(err)
	}

	percentOut := (totalOut / totalCust) * 100.0

	statList = append(statList, fmt.Sprintf("CMP: Total Customers: %s | Customers Without Power: %s (%.03f%%) | Last updated: %s.", stats.Total, stats.WithoutPower, percentOut, stats.LastUpdate.Format("2006-Jan-02/15:04")))

	for x, y := range stats.Counties {

		totalCust, err = strconv.ParseFloat(strings.ReplaceAll(y.Total, ",", ""), 64)
		if err != nil {
			panic(err)
		}

		totalOut, err = strconv.ParseFloat(strings.ReplaceAll(y.WithoutPower, ",", ""), 64)
		if err != nil {
			panic(err)
		}

		percentOut := (totalOut / totalCust) * 100.0

		statsMsg := fmt.Sprintf("CMP[County: %s]: Total Customers: %s | Customers Without Power: %s (%.03f%%) | Last updated: %s.", x, y.Total, y.WithoutPower, percentOut, stats.LastUpdate.Format("2006-Jan-02/15:04"))
		statList = append(statList, statsMsg)
	}

	return statList
}
