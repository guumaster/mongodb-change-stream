package main

import (
	"time"
  "os"
  "fmt"
  "math/rand"

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
			ClientID:  "example-publisher",
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
      "msg": gofakeit.HackerPhrase(),
      "logLevel": gofakeit.LogLevel("general"),
    }
    str, err := json.Marshal(data)
    if err != nil {
      panic(err)
    }
    msg := message.NewMessage(watermill.NewUUID(), []byte(str))
    
    err = publisher.Publish("example.topic", msg)
    if err != nil {
      panic(err)
    }

    total++
    fmt.Printf("total %5d\n", total)
    // random wait before next batch
    timer := rand.Intn(30) + 5
    time.Sleep(time.Duration(timer) * time.Millisecond)
  }
}

