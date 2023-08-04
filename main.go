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
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/hashicorp/go.net/proxy"
)

func main() {
	// Some initial Twitter setup
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessSecret := os.Getenv("ACCESS_SECRET")

	// Some initial MQTT setup
	mqttBroker := os.Getenv("MQTT_BROKER")
	mqttClientId := os.Getenv("MQTT_CLIENT_ID")
	mqttTopic := os.Getenv("MQTT_TOPIC")

	options := mqtt.NewClientOptions().AddBroker(mqttBroker).SetClientID(mqttClientId)
	options.WriteTimeout = 20 * time.Second
	mqttClient := mqtt.NewClient(options)

	// Proxy needed for hitting the CMP site
	proxyIp := os.Getenv("PROXY_IP")

	proxyDial, err := proxy.SOCKS5("tcp", proxyIp, nil, proxy.Direct)
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

	logTmp, err := readCache("/tmp/cmpLastRun")
	if err != nil {
		logTmp.Date = time.Now().Add(time.Hour * -5)
	}

	stats, err := cmpscraper.GetStats(httpClient)
	if err != nil {
		log.Fatalln("Cannot get stats:", err)
	}
	fmt.Printf("%+v\r\n", stats)

	if stats.LastUpdate.Format("2006-Jan-02/15:04") == logTmp.Date.Format("2006-Jan-02/15:04") {
		fmt.Println("No new updates.")
		return
	}

	statList := buildMessages(stats)

	// Connect to the MQTT broker
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// We limit to one message every 5 seconds, in hopes the API gods don't get mad.
	limiter := time.Tick(time.Second * 5)
	for _, y := range statList {
		<-limiter

		jsonMsg, err := json.Marshal(&MqttMessage{
			TwitterConsumerKey:    consumerKey,
			TwitterConsumerSecret: consumerSecret,
			TwitterAccessToken:    accessToken,
			TwitterAccessSecret:   accessSecret,

			Message: y,
		})
		if err != nil {
			log.Fatal(err)
		}

		token := mqttClient.Publish(mqttTopic, 2, false, jsonMsg)
		_ = token.Wait()
		if token.Error() != nil {
			panic(err)
		}
	}

	logTmp.Date = stats.LastUpdate
	err = writeCache("/tmp/cmpLastRun", logTmp)
	if err != nil {
		panic(err)
	}

	mqttClient.Disconnect(250)
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
