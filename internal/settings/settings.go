package settings

import "os"

func DBUrl() string {
	// Empty string will use the PG* variables
	return os.Getenv("DB_URL")
}
