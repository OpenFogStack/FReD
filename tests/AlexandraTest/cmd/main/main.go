package main

import (
	"os"

	"git.tu-berlin.de/mcc-fred/fred/tests/AlexandraTest/cmd/pkg/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Use pretty output since this will be used for testing
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)
	c := client.NewAlexandraClient("172.26.4.1:10000")

	// Create a keygroup
	c.CreateKeygroup("alexandraTest", true, 0, false)

	// Put a value into it.
	c.Update("alexandraTest", "id", "data", false)
}
