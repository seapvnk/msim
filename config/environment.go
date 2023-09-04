package config

type Environment string

const (
	Local	Environment = "local"
	Server	Environment = "server"
	Test	Environment	= "test"
)