package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)
type MongoInstance struct{
	Client	*mongo.Client
	db		*mongo.Database
}

var mg MongoInstance
const dbName = "mongo-go"
const mongoUrl = "mongodb://localhost:27017"+ dbName

type Employee struct{
	ID		string	`json:"id,omitempty" bson:"_id,omitempty"`
	Name	string	`json:"name"`
	Salary	float64	`json:"salary"`
	Age		float64	`json:"age"`
}

func Connect() error {
    client, err := mongo.NewClient(options.Client().ApplyURI(mongoUrl))
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }

    Db := client.Database(dbName)
    mg = MongoInstance{
        Client: client,
        db: Db,
    }

    return nil
}

func main(){
	if err := Connect(); err != nil{
		log.Fatal(err);
	}
	app := fiber.New()
	app.Get("/employee" ,func(c *fiber.Ctx) error {
		query = bson.D{{}}

		Find(c.Context(),query)
		var employees []Employee = make([Employee, 0])

	} )
}

// sfksdfdfkdf