package main

import (
	"time"
	"github.com/ibomoko/twitter-gpt-go/openai"
	"github.com/ibomoko/twitter-gpt-go/twitter"
)

func  main() {
	go scheduler()

	select{}
}


func scheduler() {
	tweetAIContent()

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tweetAIContent()
		}
	}
}

func tweetAIContent() {
	res := openai.CompleteChat("Generate an interesting fun fact about programming")
	content := res.Choices[0].Message.Content
	twitter.Tweet(content)
}