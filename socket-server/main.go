package main

import (
  "fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"net/http"
  "context"
  "log"
  "os"
  "encoding/json"

  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)


func main() {
	r := gin.Default()
	m := melody.New()

  DEBUG_MSGS := os.Getenv("DEBUG_MSGS") == "1"

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "asset/index.html")
	})

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

  client, err := connectDB()
  if err != nil {
    log.Fatal(err)
  }

  ctx := context.Background()

  collection := client.Database("demo").Collection("logs")

  pipeline := mongo.Pipeline{}
  streamOptions := options.ChangeStream().SetFullDocument(options.UpdateLookup)

  stream, err := collection.Watch(ctx, pipeline, streamOptions)
  if err != nil {
    log.Fatal(err)
  }
  log.Print("waiting for changes")

  var changeDoc map[string]interface{}

  go func() {
    r.Run("0.0.0.0:5000")
  }()

  for stream.Next(ctx) {
    if e := stream.Decode(&changeDoc); e != nil {
      log.Printf("error decoding: %s", e)
    }

    // Only show insert actions
    if changeDoc["operationType"] != "insert" {
      continue
    }

    doc := changeDoc["fullDocument"].(map[string]interface{})

    lvl := formatLevel(doc["logLevel"].(string))

    msg := fmt.Sprintf("%7s: %s", lvl, doc["msg"])

    if DEBUG_MSGS {
      log.Println(msg)
    }
		m.Broadcast([]byte(msg))

    // Uncomment this if you want to see all metadata
    //log.Printf("change: %+v", changeDoc)

    // Uncomment this to show document as JSON
    //printJSON(changeDoc["fullDocument"])
  }

}

var LEVELS = map[string]string{
  "error": "ERR",
  "warning": "WARN",
  "fatal": "FATAL",
  "debug": "DEBUG",
  "trace": "TRACE",
  "info": "INFO",
}

func formatLevel(key string) string {
  if lvl, ok := LEVELS[key]; ok {
    return fmt.Sprintf("[%s]", lvl)
  }
  return ""
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


// printJSON prints v as JSON encoded with indent to stdout. It panics on any error.
func printJSON(v interface{}) {
	w := json.NewEncoder(os.Stdout)
	w.SetIndent("", "\t")
	err := w.Encode(v)
	if err != nil {
		panic(err)
	}
}


