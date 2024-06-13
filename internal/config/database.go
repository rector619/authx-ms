package config

type Database struct {
	ConnectionString string
	DBName           string
	Migrate          bool
}
