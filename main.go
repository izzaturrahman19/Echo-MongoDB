package main

import (
	"net/http"
	"context"
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

var ctx = context.Background()

type (
	villager struct {
		Name     string `json:"name"`
		Postal   int    `json:"postal"`
	    Address  string `json:"address"`
	}
)

var (
	users = map[string]*villager{}
)


//----------
// Connect DB
//----------
func connect() (*mongo.Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client.Database("crud_echo"), nil
}

//----------
// Handlers
//----------

func createUser(c echo.Context) error {
	db, err := connect()
		if err != nil {
			log.Fatal(err.Error())
		}
	u := &villager{}
	if err := c.Bind(u); err != nil {
		return err
	}
	name := u.Name
	postal:= u.Postal
	address := u.Address

	_, err = db.Collection("villager").InsertOne(ctx, villager{name, postal, address})
		if err != nil {
			log.Fatal(err.Error())
		}
	return c.JSON(http.StatusCreated, u)
}

func getUser(c echo.Context) error {
	    db, err := connect()
		if err != nil {
			log.Fatal(err.Error())
		}
		vector := c.Param("name")

		csr, err := db.Collection("villager").Find(ctx, bson.M{"name": vector})
		if err != nil {
			log.Fatal(err.Error())
		}
		defer csr.Close(ctx)

		result := make([]villager, 0)
		for csr.Next(ctx) {
			var row villager
			err := csr.Decode(&row)
			if err != nil {
				log.Fatal(err.Error())
			}

			result = append(result, row)
		}
		
	return c.JSON(http.StatusOK, result[0])
}

func updateUser(c echo.Context) error {
		db, err := connect()
			if err != nil {
				log.Fatal(err.Error())
			}

			vector := c.Param("name")
			var selector = bson.M{"name": vector}

		csr, err := db.Collection("villager").Find(ctx, bson.M{"name": vector})
		if err != nil {
			log.Fatal(err.Error())
		}
		defer csr.Close(ctx)
	
		result := make([]villager, 0)
		for csr.Next(ctx) {
			var row villager
			err := csr.Decode(&row)
			if err != nil {
				log.Fatal(err.Error())
			}

			result = append(result, row)
		}

		u := new(villager)
		if err := c.Bind(u); err != nil {
			return err
		}
		result[0].Name = u.Name
		name:= result[0].Name
		result[0].Postal = u.Postal
		postal:= result[0].Postal
		result[0].Address = u.Address
		address:= result[0].Address
		var changes = villager{name, postal, address}

		_, err = db.Collection("villager").UpdateOne(ctx, selector, bson.M{"$set": changes})
			if err != nil {
				log.Fatal(err.Error())
			}
		return c.JSON(http.StatusOK, result[0])
}

func deleteUser(c echo.Context) error {
	db, err := connect()
			if err != nil {
				log.Fatal(err.Error())
			}

			vector := c.Param("name")
			var selector = bson.M{"name": vector}

		_, err = db.Collection("villager").DeleteOne(ctx, selector)
				if err != nil {
					log.Fatal(err.Error())
				}
			delete(users, vector)
			return c.NoContent(http.StatusNoContent)
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/villagers", createUser)
	e.GET("/villager/:name", getUser)
	e.PUT("/villager/:name", updateUser)
	e.DELETE("/villager/:name", deleteUser)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}