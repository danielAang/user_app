package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Nome      string             `json:"nome" bson:"nome,omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	sid := strings.TrimPrefix(r.URL.Path, "/usuario/")
	oid, err := primitive.ObjectIDFromHex(sid)
	switch {
	case r.Method == "GET" && err == nil:
		getUserById(w, r, oid)
	case r.Method == "GET":
		getAllUsers(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Sorry :(")
	}
}

func getDbClient() (*mongo.Client, context.Context) {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return client, ctx
}

func getUserById(w http.ResponseWriter, r *http.Request, id primitive.ObjectID) {
	c, ctx := getDbClient()
	defer c.Disconnect(ctx)
	db := c.Database("go_user_app")
	collection := db.Collection("user")
	filter := bson.M{"_id": id}
	cursor := collection.FindOne(ctx, filter)
	var user User
	err := cursor.Decode(&user)
	if err != nil {
		fmt.Fprint(w, "No data found")
		return
	}
	json, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	c, ctx := getDbClient()
	defer c.Disconnect(ctx)
	db := c.Database("go_user_app")
	collection := db.Collection("user")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Fprint(w, "No data found")
		return
	}
	var users []User
	err = cursor.All(ctx, &users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Failed to load results")
		return
	}
	json, err := json.Marshal(users)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
}
