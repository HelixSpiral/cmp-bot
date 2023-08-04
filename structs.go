package main

import "time"

type Cache struct {
	Date time.Time
}

type MqttMessage struct {
	TwitterConsumerKey    string
	TwitterConsumerSecret string
	TwitterAccessToken    string
	TwitterAccessSecret   string

	Message string
	Image   string
}
