package models

type (
	UserConfig struct {
		UserId  string
		Configs map[string]interface{}
	}
)
