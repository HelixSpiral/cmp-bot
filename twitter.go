package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dghubble/go-twitter/twitter"
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

	twitterClient := twitter.NewClient(httpClient)

	for _, message := range messages {
		log.Printf("Tweeting message: %s\r\n", message)
		tweet, resp, err := twitterClient.Statuses.Update(message, nil)
		if err != nil {
			return err
		}

		fmt.Println(tweet)
		fmt.Println(resp)
	}

	return nil
}
