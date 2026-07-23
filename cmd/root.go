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

func middlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		// lookup happens here, at command execution time
		user, err := s.db.GetUserName(context.Background(), s.cfg.UserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
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
	err = s.cfg.SetUser(cmd.args[0])
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

	err = s.cfg.SetUser(cmd.args[0])
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

	for _, user := range users {
		uStr := fmt.Sprintf("- %s", user)
		if s.cfg.UserName == user {
			uStr += " (current)"
		}
		fmt.Println(uStr)
	}

	return nil
}

// func handlerAgg(s *State, cmd Command) error {
// 	ctx := context.Background()
// 	tmpURL := "https://www.wagslane.dev/index.xml"
// 	feed, err := rss.FetchFeed(ctx, tmpURL)
// 	if err != nil {
// 		errMsg := fmt.Errorf("failed to fetch RSS feed: %v", err)
// 		return errMsg
// 	}
// 	fmt.Println(feed)
// 	return nil
// }

func handlerAddFeed(s *State, cmd Command, user database.User) error {
	ctx := context.Background()

	if len(cmd.args) != 2 {
		errMsg := errors.New("addfeed requires name, url arguments")
		return errMsg
	}

	// currUser, err := s.db.GetUserName(ctx, s.cfg.UserName)
	// if err != nil {
	// 	errMsg := fmt.Errorf("couldn't fetch currUser from db: %v", err)
	// 	return errMsg
	// }

	timeNow := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	var feedParams database.AddFeedParams
	feedParams.ID = uuid.New()
	feedParams.CreatedAt = timeNow
	feedParams.UpdatedAt = timeNow
	feedParams.UserID = user.ID
	feedParams.Name = cmd.args[0]
	feedParams.Url = cmd.args[1]

	feed, err := s.db.AddFeed(ctx, feedParams)
	if err != nil {
		errMsg := fmt.Errorf("failed to create new feed %s: %v", cmd.args[0], err)
		return errMsg
	}

	fmt.Println(feed)

	var followCmd Command
	followArgs := make([]string, 1)
	followArgs[0] = cmd.args[1]
	followCmd.name = "follow"
	followCmd.args = followArgs
	err = middlewareLoggedIn(handlerFollow)(s, followCmd)
	if err != nil {
		return err
	}

	return nil
}

func handlerFeeds(s *State, cmd Command) error {
	ctx := context.Background()

	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		errMsg := fmt.Errorf("failed to retrieve feeds from db: %v", err)
		return errMsg
	}

	for _, feed := range feeds {
		fmt.Println(feed)
	}

	return nil
}

func handlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.args) != 1 {
		errMsg := errors.New("follow cmd requires 1 url")
		return errMsg
	}

	ctx := context.Background()
	// currUser, err := s.db.GetUserName(ctx, s.cfg.UserName)
	// if err != nil {
	// 	errMsg := fmt.Errorf("couldn't fetch currUser from db: %v", err)
	// 	return errMsg
	// }
	usrFeeds, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		errMsg := fmt.Errorf("failed to get user %s's feeds: %v", s.cfg.UserName)
		return errMsg
	}

	feedExists, err := s.db.FeedExists(ctx, cmd.args[0])
	if err != nil {
		errMsg := fmt.Errorf("failed to check if feed %s exists: %v", cmd.args[0], err)
		return errMsg
	}
	if !feedExists {
		errMsg := fmt.Errorf("feed %s doesn't exists in database", cmd.args[0])
		return errMsg
	}

	feedRow, err := s.db.GetFeedByURL(ctx, cmd.args[0])
	if err != nil {
		errMsg := fmt.Errorf("failed to get feed %s: %v", cmd.args[0], err)
		return errMsg
	}

	for _, feed := range usrFeeds {
		if feed.FeedID == feedRow.ID {
			errMsg := fmt.Errorf("user %s already following feed %s", user.Name, cmd.args[0])
			return errMsg
		}
	}

	timeNow := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	var feedFollowParams database.CreateFeedFollowParams
	feedFollowParams.ID = uuid.New()
	feedFollowParams.CreatedAt = timeNow
	feedFollowParams.UpdatedAt = timeNow
	feedFollowParams.UserID = user.ID
	feedFollowParams.FeedID = feedRow.ID

	feedFollow, err := s.db.CreateFeedFollow(ctx, feedFollowParams)
	if err != nil {
		errMsg := fmt.Errorf("failed to follow feed %s: %v", cmd.args[0], err)
		return errMsg
	}
	successMsg := fmt.Sprintf("user %s following feed: %s", feedFollow.UserName, feedFollow.FeedName)
	fmt.Println(successMsg)

	return nil
}

func handlerFollowing(s *State, cmd Command, user database.User) error {
	ctx := context.Background()

	// currUser, err := s.db.GetUserName(ctx, cfg.UserName)
	// if err != nil {
	// 	errMsg := fmt.Errorf("couldn't fetch currUser from db: %v", err)
	// 	return errMsg
	// }
	usrFeeds, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		errMsg := fmt.Errorf("failed to get user %s's feeds: %v", user.Name)
		return errMsg
	}

	for _, feed := range usrFeeds {
		msg := fmt.Sprintf("feed: %s", feed.FeedName)
		fmt.Println(msg)
	}

	return nil
}

func InitCmds() (commands, error) {
	cmds := NewCommands()
	cmds.register("reset", handlerReset)
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("users", handlerUsers)
	// cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
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
