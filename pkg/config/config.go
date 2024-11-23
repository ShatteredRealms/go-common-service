package config

type ServerMode string

const (
	ModeProduction  ServerMode = "prod"
	ModeDevelopment ServerMode = "dev"
	ModeLocal       ServerMode = "local"
)

