package main

import (
	"log"

	//"github.com/Faheem-Nizar/go-rabbitmq-tutorial/internal/cmd"
	"github.com/Academics-and-Career-Council/Stargazer.git/internal/cmd"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	cmd.Execute()
}
