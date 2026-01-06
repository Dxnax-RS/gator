package commands

import(
	"errors"
	"github.com/Dxnax-RS/gator/internal/config"
)

type commands struct{
	Commands map[string]func(*config.State, command) error
}

func NewCommands() commands{
	var cmd commands
	return cmd
}

func (c *commands) Run(s *config.State, cmd command) error{
	userCommand, ok := c.Commands[cmd.Name]

	if !ok {
		return errors.New("The given command is not part of the command list")
	}

	err := userCommand(s, cmd)
	return err
}

func (c *commands) Register(name string, f func(*config.State, command) error){
	if c.Commands == nil{
		c.Commands = make(map[string]func(*config.State, command) error)
	}
	c.Commands[name] = f
}