package main

import _ "github.com/lib/pq"

import(
	"fmt"
	"errors"
	"os"
	"database/sql"
	"github.com/Dxnax-RS/gator/internal/commands"
	"github.com/Dxnax-RS/gator/internal/config"
	"github.com/Dxnax-RS/gator/internal/database"
)

func main(){

	args := os.Args
	//fmt.Println(args)

	if len(args) < 2 {
		fmt.Println(errors.New("To few arguments"))
		os.Exit(1)
	}
	cmd := commands.NewCommand()
	cmd.Name = args[1]

	if len(args) > 2 {
		cmd.Args = args[2:]
	}

	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.Db_url)
	dbQueries := database.New(db)

	var myState config.State
	myState.Db = dbQueries
	myState.Cfg = &cfg

	commandList := commands.NewCommands()
	commandList.Register("login", commands.HandlerLogin)
	commandList.Register("register", commands.RegisterUser)
	commandList.Register("reset", commands.ResetUserTable)
	commandList.Register("users", commands.GetAllUsers)
	commandList.Register("agg", commands.Aggregator)
	commandList.Register("addfeed", commands.MiddlewareLoggedIn(commands.RegisterFeed))
	commandList.Register("feeds", commands.GetAllFeeds)
	commandList.Register("follow", commands.MiddlewareLoggedIn(commands.RegisterFollow))
	commandList.Register("following", commands.MiddlewareLoggedIn(commands.GetUserFollows))
	commandList.Register("unfollow", commands.MiddlewareLoggedIn(commands.DeleteFollow))
	commandList.Register("browse", commands.MiddlewareLoggedIn(commands.BrowsePosts))
	err = commandList.Run(&myState, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}