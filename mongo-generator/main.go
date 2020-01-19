package main

import (
  "github.com/brianvoe/gofakeit/v4"
  "context"
  "os"
  "log"
  "fmt"
  "math/rand"
  "time"

  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
  maxBulkMsgs := 20

  client, err := connectDB()
  if err != nil {
    log.Fatal(err)
  }

  collection := client.Database("demo").Collection("logs")

  gofakeit.Seed(time.Now().UnixNano())

  total := 0
  for {

    // create the slice of write models
    var writes []mongo.WriteModel
    
    totalWrites := rand.Intn(maxBulkMsgs)+1

    for i := 0; i<totalWrites; i++ {
      model := mongo.NewInsertOneModel().SetDocument(bson.M{
        "msg": gofakeit.HackerPhrase(),
        "logLevel": gofakeit.LogLevel("general"),
      })
      writes = append(writes, model)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // run bulk write
    _, err := collection.BulkWrite(ctx, writes)
    if err != nil {
      log.Fatal(err)
    }

    total = total + totalWrites
    fmt.Printf("Inserted %3d Total: %6d\n", totalWrites, total)

    // random wait before next batch
    timer := rand.Intn(30) + 5
    time.Sleep(time.Duration(timer) * time.Millisecond)
  }
}

func connectDB() (*mongo.Client, error) {
  mongoURI := os.Getenv("MONGO_URI")
  if mongoURI == "" {
    mongoURI = "mongodb://localhost:30100,localhost:30101/?replicaSet=rs0&connect=direct"
  }
  fmt.Printf("Connecting... %s\n", mongoURI)
  ctx := context.TODO()
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

  fmt.Println("Connected!")

  return client, nil
}

