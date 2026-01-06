package commands

import(
	"errors"
	"fmt"
	"github.com/Dxnax-RS/gator/internal/config"
)

func HandlerLogin(s *config.State, cmd command) error{
	if len(cmd.Args) < 1{
		return errors.New("Username missing in arguments")
	}

	var err error

	err = s.Current_State.SetUser(cmd.Args[0]) 

	if err == nil{
		fmt.Println("User name updated correctly")
	}

	return err
}