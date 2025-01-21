package main

import "os"

var (
	OpenAIAPIKEY string
)

func init() {
	OpenAIAPIKEY = os.Getenv("OPENAI_API_KEY")
}
