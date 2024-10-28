package env

import (
	"os"

	"github.com/joho/godotenv"
)

var _ = godotenv.Load(".env")

var (
	LOCAL_PORT = os.Getenv("LOCAL_PORT")
	MONGO_URI  = os.Getenv("MONGO_URI")
)
