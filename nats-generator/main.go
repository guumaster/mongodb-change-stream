package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"encoding/json"
	stan "github.com/nats-io/stan.go"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/brianvoe/gofakeit/v4"
)

func main() {
	gofakeit.Seed(time.Now().UnixNano())

	natsURL := os.Getenv("STAN_URL")

	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	publisher, err := nats.NewStreamingPublisher(
		nats.StreamingPublisherConfig{
			ClusterID: "test-cluster",
			ClientID:  fmt.Sprintf("publisher-%d", time.Now().UnixNano()),
			StanOptions: []stan.Option{
				stan.NatsURL(natsURL),
			},
			Marshaler: nats.GobMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		panic(err)
	}

	// Publish a message per second
	total := 0
	for {
		data := map[string]interface{}{
			"msg":      fmt.Sprintf("[%s] %s", os.Getenv("HOSTNAME"), gofakeit.HackerPhrase()),
			"logLevel": gofakeit.LogLevel("general"),
		}
		str, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		msg := message.NewMessage(watermill.NewUUID(), []byte(str))

		err = publisher.Publish("GOLANG.topic", msg)
		if err != nil {
			panic(err)
		}

		total++
		fmt.Printf("total %5d\n", total)
		// random wait before next batch
		timer := rand.Intn(500) + 2000
		time.Sleep(time.Duration(timer) * time.Millisecond)
	}
}
