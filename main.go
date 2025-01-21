package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/webercoder/go-dalle/client"
)

func usage(msg string) {
	usage := fmt.Sprintf(`Usage: %s "prompt"`, os.Args[0])
	log.Fatalf("%s\n%s\n", msg, usage)
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		usage("no args provided")
	}

	prompt := args[0]

	c := client.NewDallEClient(OpenAIAPIKEY)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.Request(ctx, client.DallERequest{
		Model:   client.ModelDallE3,
		Prompt:  prompt,
		Size:    client.Size1792x1024,
		Quality: client.QualityStandard,
		Count:   1,
	})

	if err != nil {
		log.Fatalf("error sending request: %v", err)
	}

	for i, data := range resp.Data {
		log.Printf("Result %d: %s (revised prompt: %s)", i+1, data.URL, data.RevisedPrompt)
	}
}
