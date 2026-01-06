package config

import(
	"os"
	"encoding/json"
)

const configFileName = "/.gatorconfig.json"

type config struct{
	Db_url 				string
	Current_user_name 	string
}

func (c config) SetUser(name string) error{
	c.Current_user_name = name

	err := write(c)

	if err != nil{
		return err
	}
	
	return nil
}

func Read() (config, error){
	var configJSON config

	path, err := getConfigFilePath()

	if err != nil{
		return configJSON, err
	}

	data, err := os.ReadFile(path)

	if err != nil{
		return configJSON, err
	}

	err = json.Unmarshal(data, &configJSON)
	
	if err != nil {
		return configJSON, err
	}

	return configJSON, nil
}

func getConfigFilePath() (string, error){
	path, err := os.UserHomeDir()

	if err != nil{
		return "", err
	}

	path = path + configFileName

	return path, err
}

func write(cfg config) error{
	path, err := getConfigFilePath()

	if err != nil{
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "	")

	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0644)

	if err != nil {
		return err
	}

	return nil
}