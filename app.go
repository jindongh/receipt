package main

import (
	"embed"
	"time"
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:embed templates/*
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "templates/*"))

type Receipt struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	Text      string             `bson:"text"`
}

func getCollection(ctx context.Context) *mongo.Collection {
	dbUrl := os.Getenv("DB_URL")
	client, err := mongo.NewClient(options.Client().ApplyURI(dbUrl))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	//check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	receiptDatabase := client.Database("receipt")
	receiptsCollection := receiptDatabase.Collection("receipts")
	return receiptsCollection 
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"

	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		text := r.URL.Query().Get("text")
		collection := getCollection(ctx)
		_, err := collection.InsertOne(ctx, &Receipt{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			Text:      text,
		})
		log.Fatal(err)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		collection := getCollection(ctx)
		filter := bson.D{{}}
		cur, err := collection.Find(ctx, filter)
		if err != nil {
			log.Fatal("jindongh")
			log.Fatal(err)
		}
		receipts := []Receipt{}
		for cur.Next(ctx) {
			r := Receipt{}
			_ = cur.Decode(&t)
			receipts = append(receipts, r)
		}

		t.ExecuteTemplate(w, "index.html.tmpl", receipts)
	})

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
