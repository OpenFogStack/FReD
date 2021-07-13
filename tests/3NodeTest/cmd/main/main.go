package main

import (
	"bufio"
	"flag"
	"os"
	"strconv"
	"strings"
	"time"

	"git.tu-berlin.de/mcc-fred/fred/tests/3NodeTest/pkg/grpcclient"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Wait for the user to press enter to continue
	waitUser *bool
	reader   = bufio.NewReader(os.Stdin)
)

func main() {
	// Logging Setup
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	// Parse Flags
	waitUser = flag.Bool("wait-user", false, "wait for user input after each test")

	nodeAhost := flag.String("nodeAhost", "", "host of nodeA (e.g. localhost)") // Docker: localhost
	nodeAhttpPort := flag.String("nodeAhttp", "", "port of nodeA (e.g. 9001)")  // Docker: 9002
	nodeApeeringID := flag.String("nodeApeeringID", "", "Peering Id of nodeA")

	nodeBhost := flag.String("nodeBhost", "", "host of nodeB (e.g. localhost)")
	nodeBhttpPort := flag.String("nodeBhttp", "", "port of nodeB (e.g. 9001)")
	nodeBpeeringID := flag.String("nodeBpeeringID", "", "Peering Id of nodeB")

	nodeChost := flag.String("nodeChost", "", "host of nodeC (e.g. localhost)")
	nodeChttpPort := flag.String("nodeChttp", "", "port of nodeC (e.g. 9001)")
	nodeCpeeringID := flag.String("nodeCpeeringID", "", "Peering Id of nodeC")

	triggerNodeHost := flag.String("triggerNodeHost", "", "host of trigger node (e.g. localhost:3333)")
	triggerNodeWSHost := flag.String("triggerNodeWSHost", "", "host of trigger node web server (e.g. localhost:80)")
	triggerNodeID := flag.String("triggerNodeID", "", "Id of trigger node")

	certFile := flag.String("cert-file", "", "Certificate to talk to FReD")
	keyFile := flag.String("key-file", "", "Keyfile to talk to FReD")
	caFile := flag.String("ca-file", "", "Root certificate used to sign client certificates")

	littleCertFile := flag.String("little-cert-file", "", "Certificate to talk to FReD as \"littleclient\"")
	littleKeyFile := flag.String("little-key-file", "", "Keyfile to talk to FReD as \"littleclient\"")

	testRange := flag.String("test-range", "-", "Give tests to execute as a dash-separated range. Omitted start or end become the lowest or highest possible index, respectively. Default: All tests (\"-\"). Examples: 2-7, 1-, -6, -")

	flag.Parse()

	port, _ := strconv.Atoi(*nodeAhttpPort)
	nodeA := grpcclient.NewNode(*nodeAhost, port, *nodeApeeringID, *certFile, *keyFile, *caFile)
	port, _ = strconv.Atoi(*nodeBhttpPort)
	nodeB := grpcclient.NewNode(*nodeBhost, port, *nodeBpeeringID, *certFile, *keyFile, *caFile)
	port, _ = strconv.Atoi(*nodeChttpPort)
	nodeC := grpcclient.NewNode(*nodeChost, port, *nodeCpeeringID, *certFile, *keyFile, *caFile)
	port, _ = strconv.Atoi(*nodeAhttpPort)
	littleClient := grpcclient.NewNode(*nodeAhost, port, *nodeApeeringID, *littleCertFile, *littleKeyFile, *caFile)

	time.Sleep(15 * time.Second)

	config := &Config{
		waitUser: *waitUser,

		nodeAhost:     *nodeAhost,
		nodeAhttpPort: *nodeAhttpPort,

		nodeBhost:     *nodeBhost,
		nodeBhttpPort: *nodeBhttpPort,

		nodeChost:     *nodeChost,
		nodeChttpPort: *nodeChttpPort,

		triggerNodeHost:   *triggerNodeHost,
		triggerNodeWSHost: *triggerNodeWSHost,
		triggerNodeID:     *triggerNodeID,

		certFile: *certFile,
		keyFile:  *keyFile,

		littleCertFile: *littleCertFile,
		littleKeyFile:  *littleKeyFile,

		nodeA:        nodeA,
		nodeB:        nodeB,
		nodeC:        nodeC,
		littleClient: littleClient,
	}

	// to add a test suite, increase the size by one and add the instance of the suite to the slice
	testSuites := make([]TestSuite, 8)

	// initiate test suites
	testSuites[0] = NewStandardSuite(config)
	testSuites[1] = NewReplicaSuite(config)
	testSuites[2] = NewTriggerSuite(config)
	testSuites[3] = NewImmutableSuite(config)
	testSuites[4] = NewExpirySuite(config)
	testSuites[5] = NewSelfReplicaSuite(config)
	testSuites[6] = NewAuthenticationSuite(config)
	testSuites[7] = NewConcurrencySuite(config)

	// parse testRange, starts at 1
	minTest := 1
	maxTest := len(testSuites)
	testRangeSplit := strings.Split(*testRange, "-")
	if len(testRangeSplit[0]) > 0 {
		minTestInput, errMin := strconv.Atoi(testRangeSplit[0])
		if errMin == nil {
			if minTestInput > minTest {
				minTest = minTestInput
			}
		}
	}
	if len(testRangeSplit[1]) > 0 {
		maxTestInput, errMax := strconv.Atoi(testRangeSplit[1])
		if errMax == nil {
			if maxTestInput < maxTest {
				maxTest = maxTestInput
			}
		}
	}
	if minTest > maxTest {
		minTest = maxTest
	}

	// run tests
	// minTest and maxTest are indexed with 1 at the beginning, but the slice starts at 0
	for i := minTest - 1; i < maxTest; i++ {
		testSuites[i].RunTests()
	}

	// tally errors
	totalerrors := nodeA.Errors + nodeB.Errors + nodeC.Errors + littleClient.Errors

	if totalerrors > 0 {
		log.Error().Msgf("Total Errors: %d", totalerrors)
	}

	os.Exit(totalerrors)
}

func logNodeAction(node *grpcclient.Node, format string, a ...interface{}) {
	wait()
	log.Info().Str("node", node.ID).Msgf(format, a...)
}

func logNodeFailure(node *grpcclient.Node, expected, result string) {
	wait()
	log.Warn().Str("node", node.ID).Msgf("expected: %s, but got: %#v", expected, result)
	node.Errors++
}

func wait() {
	if *waitUser {
		log.Info().Msg("Please press enter to continue:")
		_, _, _ = reader.ReadLine()
	} else {
		time.Sleep(1 * time.Second)
	}
}
