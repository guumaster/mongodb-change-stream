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
  mongoURI := os.Getenv("MONGO_URI")
  if mongoURI == "" {
    mongoURI = "mongodb://localhost:30100,localhost:30101/?replicaSet=rs0&connect=direct"
  }
  fmt.Printf("Connecting... %s\n", mongoURI)
  ctx := context.TODO()
  clientOptions := options.Client().ApplyURI(mongoURI)
  client, err := mongo.Connect(ctx, clientOptions)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Ping...")
  err = client.Ping(ctx, nil)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("Connected!")

  collection := client.Database("demo").Collection("logs")

  gofakeit.Seed(time.Now().UnixNano())

  total := 0
  for {
    level := gofakeit.LogLevel("general")
    msg := gofakeit.HackerPhrase()
    

    // create the slice of write models
    var writes []mongo.WriteModel
    
    totalWrites := rand.Intn(10)+1

    // range over each list of operations and create the write model
    for i := 0; i<totalWrites; i++ {
      model := mongo.NewInsertOneModel().SetDocument(bson.M{
        "msg": msg,
        "logLevel": level,
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
    /*
    _, err := collection.InsertOne(ctx, bson.M{"msg": msg, "logLevel": level })
    if err != nil {
      log.Fatal(err)
    }
    */

    total = total + totalWrites
    fmt.Printf("Inserted %5d\n", total)

    timer := rand.Intn(30) + 10
    time.Sleep(time.Duration(timer) * time.Millisecond)
  }


}

