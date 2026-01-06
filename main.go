package main

import(
	"fmt"
	"errors"
	"os"
	"github.com/Dxnax-RS/gator/internal/commands"
	"github.com/Dxnax-RS/gator/internal/config"
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

	var myState config.State
	myState.Current_State = &cfg

	commandList := commands.NewCommands()
	commandList.Register("login", commands.HandlerLogin)
	err = commandList.Run(&myState, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}