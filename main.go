package main

import (
	"fmt"
	"gator/internal/config"
	"os"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}

	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	if c.handlers == nil {
		c.handlers = make(map[string]func(*state, command) error)
	}
	c.handlers[name] = f
}

func handlerLogin(s *state, cmd command) error {
	cfg, configCreated, err := config.ReadOrCreate()
	if err != nil {
		return err
	}
	s.cfg = &cfg

	if len(cmd.args) == 0 {
		if configCreated {
			fmt.Println("Config file has been created")
			return nil
		}
		return fmt.Errorf("username is required")
	}

	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return err
	}

	fmt.Printf("User has been set to %s\n", cmd.args[0])
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "error: not enough arguments provided")
		os.Exit(1)
	}

	appState := state{}

	cmds := commands{
		handlers: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	if err := cmds.run(&appState, cmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
