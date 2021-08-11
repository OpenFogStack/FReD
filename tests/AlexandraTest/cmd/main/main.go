package main

import (
	"os"
	"time"

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
	time.Sleep(20 * time.Second)
	c := client.NewAlexandraClient("172.26.4.1:10000")

	// Create a keygroup
	log.Info().Msg("Creating Keygroup alexandraTest")
	c.CreateKeygroup("nodeB", "alexandraTest", true, 10_000, false)

	// Put a value into it.
	log.Info().Msg("Putting Value into alexandraTest")
	c.Update("alexandraTest", "id", "data", false)

	// Read this value from anywhere
	log.Info().Msg("Reading Value from alexandraTest")
	read := c.Read("alexandraTest", "id", 500, false)

	if len(read) != 1 || read[0] != "data" {
		log.Fatal().Msgf("Read alexandraTest/id: expected 'data' but got '%v'", read)
	}

	// Add the other nodes to the keygroup
	log.Info().Msg("Adding nodeA and nodeC as replicas")
	c.AddKeygroupReplica("alexandraTest", "nodeA", 600, false)
	c.AddKeygroupReplica("alexandraTest", "nodeC", 600, false)

	log.Info().Msg("Reading Value from alexandraTest again")
	read = c.Read("alexandraTest", "id", 500, false)

	if len(read) != 1 || read[0] != "data" {
		log.Fatal().Msgf("Read alexandraTest/id: expected 'data' but got '%v'", read)
	}
	log.Info().Msgf("Finished!")
}
