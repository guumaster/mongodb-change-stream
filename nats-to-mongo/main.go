package main

import (
  "context"
  "log"
	"time"
  "os"
  "fmt"
  "encoding/json"

	stan "github.com/nats-io/stan.go"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"

  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo/options"
)


func main() {
  client, err := connectDB()
  if err != nil {
    log.Fatal(err)
  }

  subscriber, err := connectNats()
  if err != nil {
    log.Fatal(err)
  }

  ctxNats := context.Background()
  collection := client.Database("demo").Collection("logs")

	messages, err := subscriber.Subscribe(ctxNats, "example.topic")
	if err != nil {
    log.Fatal(err)
	}

  total := 0
  for msg := range messages {
    total++
    log.Printf("total: %6d", total)

    var bsonMap bson.M
    err := json.Unmarshal(msg.Payload, &bsonMap)
    if err != nil {
      log.Fatal(err)
    }
    //printJSON(bsonMap)

    _, err = collection.InsertOne(context.Background(), bsonMap)
    if err != nil {
      msg.Nack()
      log.Fatal("INSERT FAIL ",  err)
    } else {
      // we need to Acknowledge that we received and processed the message,
      // otherwise, it will be resent over and over again.
      msg.Ack()
    }
  }
}

func connectNats() (*nats.StreamingSubscriber, error) {
  natsURL := os.Getenv("STAN_URL")

  if natsURL == "" {
    natsURL = "nats://localhost:4222"
  }

  subscriber, err := nats.NewStreamingSubscriber(
    nats.StreamingSubscriberConfig{
      ClusterID:        "test-cluster",
      ClientID:         "example-subscriber",
      QueueGroup:       "example",
      DurableName:      "my-durable",
      SubscribersCount: 4, // how many goroutines should consume messages
      //CloseTimeout:     time.Minute,
      AckWaitTimeout:   time.Second * 30,
      StanOptions: []stan.Option{
        stan.NatsURL(natsURL),
      },
      Unmarshaler: nats.GobMarshaler{},
    },
    watermill.NewStdLogger(false, false),
  )
  if err != nil {
    return nil, err
  }
  fmt.Println("NATS connected!")

  return subscriber, nil
}

func connectDB() (*mongo.Client, error) {
  mongoURI := os.Getenv("MONGO_URI")
  if mongoURI == "" {
    mongoURI = "mongodb://localhost:30100,localhost:30101/?replicaSet=rs0&connect=direct"
  }
  fmt.Printf("Connecting... %s\n", mongoURI)
  ctx := context.Background()
  clientOptions := options.Client().ApplyURI(mongoURI)
  client, err := mongo.Connect(ctx, clientOptions)
  if err != nil {
    return nil, err
  }
  fmt.Println("Ping...")
  err = client.Ping(ctx, nil)
  if err != nil {
    return nil, err
  }

  fmt.Println("Mongo connected!")

  return client, nil
}


// printJSON prints v as JSON encoded with indent to stdout. It panics on any error.
func printJSON(v interface{}) {
	w := json.NewEncoder(os.Stdout)
	w.SetIndent("", "\t")
	err := w.Encode(v)
	if err != nil {
		panic(err)
	}
}

