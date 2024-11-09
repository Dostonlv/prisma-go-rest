package main

import (
	"context"
	"log"

	"github.com/Dostonlv/prisma-go-rest/db"
	"github.com/gofiber/fiber/v2"
)

var Client = db.NewClient(
	db.WithDatasourceURL("postgresql://delta_owner:ZTz5jYshP6Og@ep-mute-credit-a5phlwm0.us-east-2.aws.neon.tech/delta?schema=public"),
)

func main() {

	if err := Client.Prisma.Connect(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := Client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()
	app := fiber.New()


	app.Get("/users", getUsers)
	app.Get("/users/:id", getUser)
	app.Post("/users", createUser)
	app.Put("/users/:id", updateUser)
	app.Delete("/users/:id", deleteUser)

	log.Fatal(app.Listen(":3000"))
}

func getUsers(c *fiber.Ctx) error {
	users, err := Client.User.FindMany().Exec(context.Background())
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(users)
}

func getUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := Client.User.FindUnique(db.User.ID.Equals(id)).Exec(context.Background())
	if err != nil {
		return c.Status(404).SendString("User not found")
	}
	return c.JSON(user)
}

func createUser(c *fiber.Ctx) error {
	data := new(struct {
		FullName string `json:"fullName"`
		Email    string `json:"email"`
	})
	if err := c.BodyParser(data); err != nil {
		return c.Status(400).SendString("Invalid request body")
	}

	user, err := Client.User.CreateOne(
		db.User.FullName.Set(data.FullName),
		db.User.Email.Set(data.Email),
	).Exec(context.Background())

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.Status(201).JSON(user)
}

func updateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	data := new(struct {
		FullName string `json:"fullName"`
		Email    string `json:"email"`
	})
	if err := c.BodyParser(data); err != nil {
		return c.Status(400).SendString("Invalid request body")
	}

	_, err := Client.User.FindUnique(db.User.ID.Equals(id)).Exec(context.Background())
	if err != nil {
		return c.Status(404).SendString("User not found")
	}

	updatedUser, err := Client.User.FindUnique(db.User.ID.Equals(id)).Update(
		db.User.FullName.Set(data.FullName),
		db.User.Email.Set(data.Email),
	).Exec(context.Background())

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(updatedUser)
}

func deleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := Client.User.FindUnique(db.User.ID.Equals(id)).Delete().Exec(context.Background())
	if err != nil {
		return c.Status(404).SendString("User not found")
	}
	return c.SendStatus(204)
}

// func run() error {

// 	if err := Client.Prisma.Connect(); err != nil {
// 		return err
// 	}
// 	defer func() {
// 		if err := Client.Prisma.Disconnect(); err != nil {
// 			panic(err)
// 		}
// 	}()
// 	ctx := context.Background()
// 	createUser, err := Client.User.CreateOne(
// 		db.User.FullName.Set("Apple"),
// 		db.User.Email.Set("apple@apple.com")).Exec(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	result, err := json.MarshalIndent(createUser, "", " ")
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Printf("createdUser: %s\n", result)

// 	user, err := Client.User.FindUnique(
// 		db.User.ID.Equals(createUser.ID),
// 	).Exec(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	result, err = json.MarshalIndent(user, "", " ")
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Printf("User: %s\n", result)

// 	return nil
// }
