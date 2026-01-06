package commands

type command struct{
	Name string
	Args []string
}

func NewCommand() command{
	var cmd command
	return cmd
}