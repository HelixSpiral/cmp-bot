package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/HelixSpiral/cmpscraper"
	"github.com/hashicorp/go.net/proxy"
)

func main() {
	//proxyDial, err := proxy.SOCKS5("tcp", "72.210.221.197:4145", nil, proxy.Direct)
	proxyDial, err := proxy.SOCKS5("tcp", "68.71.247.130:4145", nil, proxy.Direct)
	if err != nil {
		log.Fatalln("Cannot connect to proxy:", err)
	}

	httpTransport := &http.Transport{
		Dial: proxyDial.Dial,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	httpClient := &http.Client{
		Transport: httpTransport,
	}

	stats, err := cmpscraper.GetStats(httpClient)
	if err != nil {
		log.Fatalln("Cannot get stats:", err)
	}
	fmt.Printf("%+v\r\n", stats)

	logTmp, err := readCache("/tmp/cmpLastRun")
	if err != nil {
		logTmp.Date = time.Now().Add(time.Hour * -5)
	}

	if stats.LastUpdate == logTmp.Date {
		fmt.Println("No new updates.")
		return
	}

	statList := buildMessages(stats)

	err = sendMessages(statList)
	if err != nil {
		panic(err)
	}

	logTmp.Date = stats.LastUpdate
	err = writeCache("/tmp/cmpLastRun", logTmp)
	if err != nil {
		panic(err)
	}
}

func readCache(f string) (Cache, error) {
	var logTmp Cache
	rawdata, err := os.ReadFile(f)
	if err != nil {
		return Cache{}, err
	}

	err = json.Unmarshal(rawdata, &logTmp)
	if err != nil {
		return Cache{}, err
	}

	return logTmp, nil
}

func writeCache(f string, logTmp Cache) error {
	jsonData, err := json.Marshal(logTmp)
	if err != nil {
		return err
	}

	_ = jsonData
	err = os.WriteFile(f, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
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
