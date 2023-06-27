package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dghubble/oauth1"
)

func sendMessages(messages []string) error {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessSecret := os.Getenv("ACCESS_SECRET")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	for _, message := range messages {
		log.Printf("Tweeting message: %s\r\n", message)

		resp, err := httpClient.Post("https://api.twitter.com/2/tweets", "application/json",
			bytes.NewBuffer([]byte(fmt.Sprintf(`{"text": "%s"}`, message))))
		if err != nil {
			return err
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		log.Println("Tweet:", string(body))
	}

	return nil
}
