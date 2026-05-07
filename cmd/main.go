package main

import (
	"log"

	"github.com/chopic82region/tz-junior.git/internal/config"
)

func main() {

	_, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error of load config")
	}

}
