package main

import (
	"bufio"
	"fmt"
	"github.com/zeromq/goczmq"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/zmqclient"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/zmqserver"
	"log"
	"math"
	"os"
	"time"
)

func main() {

	multipart := true

	rec, _ := zmqserver.NewReceiver("test", "5555")

	var start int64 = 0
	var end int64 = 0

	go func() {
		for {
			recvd, err := rec.Receive()
			end = time.Now().UnixNano()
			log.Printf("Receiver: got message with size: %d with error %v", len(recvd), err)
			log.Printf("Took me: %d milliseconds", (end-start)/1000000)
		}
	}()

	// Create a very big string
	oneGb := int(math.Pow(10, 9))
	b := make([]byte, (oneGb*30)/10)

	for i := range b {
		b[i] = 0x41 // ASCII "A"
	}

	if multipart {
		dealer, _ := goczmq.NewDealer(fmt.Sprintf("tcp://%s:%d", "localhost", 5555))
		// So that the router doesnt send an answer
		b = append([]byte{0x00}, b...)
		start = time.Now().UnixNano()
		parts := 60
		for i := 0; i < len(b); i += len(b) / parts { // Send in 10 parts
			end := i + len(b)/parts

			flag := goczmq.FlagMore

			if end >= len(b) {
				end = len(b)
				flag = goczmq.FlagNone
				log.Printf("Sending last frame")
			}

			bytes := b[i:end]
			log.Printf("Sending frame #%d", i)
			dealer.SendFrame(bytes, flag)
		}
		b = nil
		//log.Printf("Sending last frame")
		//dealer.SendFrame([]byte{0x41}, goczmq.FlagNone) // Send one byte more, why not
		log.Printf("Done sending! Press enter to continue:")
		bufio.NewReader(os.Stdin).ReadString('\n')
	} else {
		log.Printf("Converting byte[] to string")
		log.Printf("Created a very big string of size %dGB", len(b)/int(math.Pow(10, 9)))
		// Send a very big file in one message
		sen := zmqclient.NewSender("localhost", 5555)
		log.Printf("Sending and waiting for answer")
		start = time.Now().UnixNano()
		_ = sen.SendBytes(b)
		b = nil
		log.Printf("Done sending! Press enter to continue:")
		bufio.NewReader(os.Stdin).ReadString('\n')
	}
}
