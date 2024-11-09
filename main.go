package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Dostonlv/prisma-go-rest/db"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	client := db.NewClient(
		db.WithDatasourceURL("postgresql://delta_owner:ZTz5jYshP6Og@ep-mute-credit-a5phlwm0.us-east-2.aws.neon.tech/delta?schema=public"),
	)
	if err := client.Prisma.Connect(); err != nil {
		return err
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()
	ctx := context.Background()
	createUser, err := client.User.CreateOne(
		db.User.FullName.Set("Apple"),
		db.User.Email.Set("apple@apple.com")).Exec(ctx)
	if err != nil {
		return err
	}

	result, err := json.MarshalIndent(createUser, "", " ")
	if err != nil {
		return err
	}
	fmt.Printf("createdUser: %s\n", result)

	user, err := client.User.FindUnique(
		db.User.ID.Equals(createUser.ID),
	).Exec(ctx)
	if err != nil {
		return err
	}
	result, err = json.MarshalIndent(user, "", " ")
	if err != nil {
		return err
	}
	fmt.Printf("User: %s\n", result)

	return nil
}
