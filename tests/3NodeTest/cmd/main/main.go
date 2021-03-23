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
	waitUser bool
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
	waitUser = *flag.Bool("wait-user", false, "wait for user input after each test")

	nodeAhost := *flag.String("nodeAhost", "172.26.1.1", "host of nodeA (e.g. localhost)") // Docker: localhost
	nodeAhttpPort := *flag.String("nodeAhttp", "9001", "port of nodeA (e.g. 9001)")        // Docker: 9002
	nodeApeeringID := *flag.String("nodeAzmqID", "nodeA", "ZMQ Id of nodeA")

	nodeBhost := *flag.String("nodeBhost", "172.26.2.1", "host of nodeB (e.g. localhost)")
	nodeBhttpPort := *flag.String("nodeBhttp", "9001", "port of nodeB (e.g. 9001)")
	nodeBpeeringID := *flag.String("nodeBzmqID", "nodeB", "ZMQ Id of nodeB")

	nodeChost := *flag.String("nodeChost", "172.26.3.1", "host of nodeC (e.g. localhost)")
	nodeChttpPort := *flag.String("nodeChttp", "9001", "port of nodeC (e.g. 9001)")
	nodeCpeeringID := *flag.String("nodeCzmqID", "nodeC", "ZMQ Id of nodeC")

	triggerNodeHost := *flag.String("triggerNodeHost", "172.26.5.1:3333", "host of trigger node (e.g. localhost:3333)")
	triggerNodeWSHost := *flag.String("triggerNodeWSHost", "172.26.5.1:80", "host of trigger node web server (e.g. localhost:80)")
	triggerNodeID := *flag.String("triggerNodeID", "triggernode", "Id of trigger node")

	certFile := *flag.String("cert-file", "/cert/client.crt", "Certificate to talk to FReD")
	keyFile := *flag.String("key-file", "/cert/client.key", "Keyfile to talk to FReD")

	littleCertFile := *flag.String("little-cert-file", "/cert/littleclient.crt", "Certificate to talk to FReD as \"littleclient\"")
	littleKeyFile := *flag.String("little-key-file", "/cert/littleclient.key", "Keyfile to talk to FReD as \"littleclient\"")

	testRange := *flag.String("test-range", "-", "Give tests to execute as a dash-separated range. Omitted start or end become the lowest or highest possible index, respectively. Default: All tests (\"-\"). Examples: 2-7, 1-, -6, -")

	flag.Parse()

	port, _ := strconv.Atoi(nodeAhttpPort)
	nodeA := grpcclient.NewNode(nodeAhost, port, certFile, keyFile)
	port, _ = strconv.Atoi(nodeBhttpPort)
	nodeB := grpcclient.NewNode(nodeBhost, port, certFile, keyFile)
	port, _ = strconv.Atoi(nodeChttpPort)
	nodeC := grpcclient.NewNode(nodeChost, port, certFile, keyFile)
	port, _ = strconv.Atoi(nodeAhttpPort)
	littleClient := grpcclient.NewNode(nodeAhost, port, littleCertFile, littleKeyFile)

	time.Sleep(15 * time.Second)

	config := &Config{
		waitUser: waitUser,

		nodeAhost:      nodeAhost,
		nodeAhttpPort:  nodeAhttpPort,
		nodeApeeringID: nodeApeeringID,

		nodeBhost:      nodeBhost,
		nodeBhttpPort:  nodeBhttpPort,
		nodeBpeeringID: nodeBpeeringID,

		nodeChost:      nodeChost,
		nodeChttpPort:  nodeChttpPort,
		nodeCpeeringID: nodeCpeeringID,

		triggerNodeHost:   triggerNodeHost,
		triggerNodeWSHost: triggerNodeWSHost,
		triggerNodeID:     triggerNodeID,

		certFile: certFile,
		keyFile:  keyFile,

		littleCertFile: littleCertFile,
		littleKeyFile:  littleKeyFile,

		nodeA: nodeA,
		nodeB: nodeB,
		nodeC: nodeC,
		littleClient: littleClient,
	}

	// to add a test suite, increase the size by one and add the instance of the suite to the slice
	testSuites := make([]TestSuite, 7)

	// initiate test suites
	testSuites[0] = NewStandardSuite(config)
	testSuites[1] = NewReplicaSuite(config)
	testSuites[2] = NewTriggerSuite(config)
	testSuites[3] = NewImmutableSuite(config)
	testSuites[4] = NewExpirySuite(config)
	testSuites[5] = NewSelfReplicaSuite(config)
	testSuites[6] = NewAuthenticationSuite(config)

	// parse testRange, starts at 1
	minTest := 1
	maxTest := len(testSuites)
	testRangeSplit := strings.Split(testRange, "-")
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
	for i := minTest-1; i < maxTest ; i++ {
		testSuites[i].RunTests()
	}

	// tally errors
	totalerrors := nodeA.Errors + nodeB.Errors + nodeC.Errors + littleClient.Errors

	if totalerrors > 0 {
		log.Error().Msgf("Total Errors: %d", totalerrors)
	}

	os.Exit(totalerrors)
}

func logNodeAction(node *grpcclient.Node, action string) {
	wait()
	log.Info().Str("node", node.Addr).Msg(action)
}

func logNodeFailure(node *grpcclient.Node, expected, result string) {
	wait()
	log.Warn().Str("node", node.Addr).Msgf("expected: %s, but got: %#v", expected, result)
	node.Errors++
}

func wait() {
	if waitUser {
		log.Info().Msg("Please press enter to continue:")
		_, _, _ = reader.ReadLine()
	} else {
		time.Sleep(1 * time.Second)
	}
}
