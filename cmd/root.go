package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/g-lok/bootdev-gatorblogs/internal/config"
)

type State struct {
	cfg *config.Config
}

type Command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*State, Command) error
}

// Constructor function ensures the map is never nil
func NewCommands() *commands {
	return &commands{
		commands: make(map[string]func(*State, Command) error),
	}
}

func (c *commands) register(name string, f func(*State, Command) error) {
	c.commands[name] = f
}

func (c *commands) run(s *State, cmd Command) error {
	cmdFunc, exists := c.commands[cmd.name]
	if !exists {
		errMsg := fmt.Errorf("command %v not found", cmd.name)
		return errMsg
	}

	err := cmdFunc(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func getConfig() (config.Config, error) {
	currentCfg, err := config.Read()
	if err != nil {
		errMsg := fmt.Errorf("failed to read .gatorconfig.json: %w", err)
		return config.Config{}, errMsg
	}

	return *currentCfg, nil
}

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.args) != 1 {
		errMsg := errors.New("login cmd requires 1 username argument")
		return errMsg
	}

	cfg, err := getConfig()
	if err != nil {
		errMsg := fmt.Errorf("failed to get ~/.gatorconfig: %v", err)
		return errMsg
	}

	err = cfg.SetUser(cmd.args[0])
	if err != nil {
		errMsg := errors.New("login cmd requires 1 username argument")
		return errMsg
	}

	fmt.Printf("username has been set to %s\n", cmd.args[0])

	return nil
}

func handlerRegister(s *State, cmd Command) error {
	return nil
}

func handlerUsers(s *State, cmd Command) error {
	return nil
}

func InitCmds() (commands, error) {
	cmds := NewCommands()
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("users", handlerUsers)
	return *cmds, nil
}

func (s *State) InitState() error {
	cfg, err := getConfig()
	if err != nil {
		errMsg := fmt.Errorf("failed to load config: %w", err)
		return errMsg
	}

	s.cfg = &cfg
	return nil
}

func Root() error {
	cmds, err := InitCmds()
	if err != nil {
		return err
	}

	var state State
	err = state.InitState()
	if err != nil {
		return err
	}

	args := os.Args
	if len(args) <= 1 {
		errMsg := errors.New("cli requires at least 1 argument")
		return errMsg
	}

	// fmt.Println(args)
	var cmd Command
	cmd.name = args[1]
	cmd.args = args[2:]

	err = cmds.run(&state, cmd)
	if err != nil {
		return err
	}

	return nil
}
