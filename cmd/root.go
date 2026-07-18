package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/g-lok/bootdev-gatorblogs/internal/config"
	"github.com/g-lok/bootdev-gatorblogs/internal/database"
	"github.com/google/uuid"
)

type State struct {
	db  *database.Queries
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

func handlerReset(s *State, cmd Command) error {
	ctx := context.Background()
	err := s.db.ResetUsers(ctx)
	if err != nil {
		errMsg := errors.New("failed to reset table 'users'")
		return errMsg
	}

	fmt.Println("table 'users' has been reset")
	return nil
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

	ctx := context.Background()
	usrExists, err := s.db.UserExists(ctx, cmd.args[0])
	if err != nil {
		errMsg := fmt.Errorf("failed to retrieve user %s: %v", cmd.args[0], err)
		return errMsg
	}

	if !usrExists {
		errMsg := fmt.Errorf("user %s does not exist in database", cmd.args[0])
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
	if len(cmd.args) != 1 {
		errMsg := errors.New("register cmd requires 1 username argument")
		return errMsg
	}

	ctx := context.Background()
	usrExists, err := s.db.UserExists(ctx, cmd.args[0])
	if err != nil {
		errMsg := fmt.Errorf("failed to retrieve user %s: %v", cmd.args[0], err)
		return errMsg
	}
	if usrExists {
		errMsg := fmt.Errorf("user %s already exists in database", cmd.args[0])
		return errMsg
	}

	timeNow := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	var userParams database.CreateUserParams
	userParams.ID = uuid.New()
	userParams.Name = cmd.args[0]
	userParams.CreatedAt = timeNow
	userParams.UpdatedAt = timeNow

	usr, err := s.db.CreateUser(ctx, userParams)
	if err != nil {
		errMsg := fmt.Errorf("error registering user %s: %v", cmd.args[0], err)
		return errMsg
	}

	cfg, err := getConfig()
	if err != nil {
		errMsg := fmt.Errorf("failed to load config: %w", err)
		return errMsg
	}
	err = cfg.SetUser(cmd.args[0])
	if err != nil {
		errMsg := fmt.Errorf("failed to set user %s: %v", cmd.args[0], err)
		return errMsg
	}

	fmt.Printf("user %s registered\n", cmd.args[0])
	log.Print(usr)

	return nil
}

func handlerUsers(s *State, cmd Command) error {
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		errMsg := fmt.Errorf("failed to retrieve users: %v", err)
		return errMsg
	}

	cfg, err := getConfig()
	if err != nil {
		errMsg := fmt.Errorf("failed to get ~/.gatorconfig: %v", err)
		return errMsg
	}

	for _, user := range users {
		s := fmt.Sprintf("- %s", user)
		if cfg.UserName == user {
			s += " (current)"
		}
		fmt.Println(s)
	}

	return nil
}

func InitCmds() (commands, error) {
	cmds := NewCommands()
	cmds.register("reset", handlerReset)
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
	db, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		errMsg := fmt.Errorf("failed to open db %s: %v", cfg.URL, err)
		return errMsg
	}
	dbQueries := database.New(db)
	// defer db.Close()

	s.cfg = &cfg
	s.db = dbQueries

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
