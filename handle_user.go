package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joshalling/gatorcli/internal/database"
)

func handleLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("the login command expects one argument: username")
	}
	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	err = s.c.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Logged in as %s\n", cmd.args[0])

	return nil
}

func handleRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("the login command expects one argument: username")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		Name:      cmd.args[0],
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	err = s.c.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Registered as %s\n", user.Name)
	fmt.Printf("User ID: %s\n", user.ID)
	fmt.Printf("Created at: %s\n", user.CreatedAt)
	fmt.Printf("Updated at: %s\n", user.UpdatedAt)

	return nil
}
