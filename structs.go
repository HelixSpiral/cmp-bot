package main

import "time"

type Cache struct {
	Date time.Time
}

type MqttMessage struct {
	MastodonServer       string
	MastodonClientID     string
	MastodonClientSecret string
	MastodonUser         string
	MastodonPass         string

	TwitterConsumerKey    string
	TwitterConsumerSecret string
	TwitterAccessToken    string
	TwitterAccessSecret   string

	Message string
	Image   string
}
