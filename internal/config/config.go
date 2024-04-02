package config

type Config struct {
	History        map[string]string `json:"history"`
	Authorizations map[string]string `json:"authorizations"`
}
