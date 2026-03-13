package config

import(
	"github.com/Dxnax-RS/gator/internal/database"
)

type State struct{
	Db  *database.Queries
	Cfg *config
}