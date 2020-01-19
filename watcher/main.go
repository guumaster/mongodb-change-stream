package main

import (
  "context"
  "log"
  "fmt"
  "os"
  "encoding/json"

  "go.mongodb.org/mongo-driver/mongo"
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

  pipeline := mongo.Pipeline{}
  streamOptions := options.ChangeStream().SetFullDocument(options.UpdateLookup)

  stream, err := collection.Watch(ctx, pipeline, streamOptions)
  if err != nil {
    log.Fatal(err)
  }
  log.Print("waiting for changes")

  var changeDoc map[string]interface{}

  for stream.Next(ctx) {
    if e := stream.Decode(&changeDoc); e != nil {
      log.Printf("error decoding: %s", e)
    }

    // Uncomment this if you want to see all metadata
    //log.Printf("change: %+v", changeDoc)

    printJSON(changeDoc["fullDocument"])
  }

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

