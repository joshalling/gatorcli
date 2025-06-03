package main

import (
	"fmt"
	"os"

	"github.com/joshalling/gatorcli/internal/config"
)

type state struct {
	c *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v", err)
		os.Exit(1)
	}

	s := &state{
		c: &c,
	}

	cmds := commands{handlers: make(map[string]func(*state, command) error)}
	cmds.register("login", handleLogin)

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Please provide a command")
		os.Exit(1)
	}

	err = cmds.run(s, command{name: args[0], args: args[1:]})
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
		os.Exit(1)
	}
}

func handleLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("the login command expects one argument: username")
	}

	err := s.c.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Logged in as %s\n", cmd.args[0])

	return nil
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.handlers[cmd.name]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}

	return handler(s, cmd)
}

func (c *commands) register(name string, handler func(*state, command) error) {
	c.handlers[name] = handler
}
