package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	db     *mongo.Database
}

var mg MongoInstance

const dbName = "mongo-go"
const mongoURI = "mongodb://localhost:27017/" + dbName

type Employee struct {
	ID     string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
	Age    float64 `json:"age"`
}

func Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
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
		db:     Db,
	}

	return nil
}

func main() {
	if err := Connect(); err != nil {
		log.Fatal(err)
	}
	app := fiber.New()
	app.Get("/employee", func(c *fiber.Ctx) error {

		query := bson.D{{}}

		cursor, err := mg.db.Collection("employees").Find(c.Context(), query)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		var employees []Employee = make([]Employee, 0)

		if err := cursor.All(c.Context(), employees); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(employees)
	})
	app.Post("/employee", func(c *fiber.Ctx) error {

		collection := mg.db.Collection("employees")
		employee := new(Employee)
		if err := c.BodyParser(employee); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		inResult, err := collection.InsertOne(c.Context(), employee)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		filter := bson.D{{Key: "_id", Value: inResult.InsertedID}}

		createdRecord := collection.FindOne(c.Context(), filter)

		createdEmployee := &Employee{}
		if err := createdRecord.Decode(createdEmployee); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.Status(201).JSON(createdEmployee)

	})
	app.Put("/employee/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		employeeId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		collection := mg.db.Collection("employees")
		employee := new(Employee)

		if err := c.BodyParser(employee); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		filter := bson.D{{Key: "_id", Value: employeeId}}
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: employee.Name},
				{Key: "age", Value: employee.Age},
				{Key: "salary", Value: employee.Salary},
			}},
		}

		updateErr := collection.FindOneAndUpdate(c.Context(), filter, update).Err()
		if updateErr != nil {
			if updateErr == mongo.ErrNoDocuments {
				return c.Status(404).SendString("Employee not found")
			}
			return c.Status(500).SendString(updateErr.Error())
		}

		employee.ID = id

		return c.Status(200).JSON(employee)
	})

	app.Delete("/employee/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		employeeId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		filter := bson.D{{
			Key:   "_id",
			Value: employeeId,
		}}
		collection := mg.db.Collection("employees")
		err = collection.FindOneAndDelete(c.Context(), filter).Err()
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(404).SendString("Employee not found")
			}
			return c.Status(500).SendString(err.Error())
		}
		return c.Status(200).SendString("Employee deleted")
	})
	log.Fatal(app.Listen(":3000"))
}
