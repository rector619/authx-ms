package config

type ServerConfiguration struct {
	Port                          string
	Secret                        string
	AccessTokenExpirationDuration int
}
type App struct {
	Name string
	Key  string
}
