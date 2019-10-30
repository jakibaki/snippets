package main

// if you cheat with this I will personally throw you out

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

func getPixString(x int, y int, color string) string {
	return fmt.Sprintf("CC %d %d %s\n", x, y, color)
}

func getPart(msg []string, part int, total int) string {
	var out strings.Builder

	chunkSize := (len(msg) + total - 1) / total
	start := chunkSize * part

	out.Grow(len("CC 0000 0000 000000") * 1920 * 1080)
	for i := start; i < start+chunkSize; i++ {
		fmt.Fprint(&out, msg[i])
	}
	return out.String()
}

func senderThread(message chan string) {
	conn, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	msg := ""

	for {
		select {
		case msg = <-message:
		default:
			if len(msg) != 0 {
				fmt.Fprint(conn, msg)
			} else {
				time.Sleep(time.Millisecond * 10)
			}
		}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	numThreads := 10
	resX := 1920
	resY := 1080
	var threadMessages []chan string

	for i := 0; i < numThreads; i++ {
		threadMessages = append(threadMessages, make(chan string))
		go senderThread(threadMessages[i])
	}

	for scanner.Scan() {
		color := scanner.Text()
		var msg = make([]string, resX*resY)

		for x := 0; x < resX; x++ {
			for y := 0; y < resY; y++ {
				msg[x+y*resX] = getPixString(x, y, color)
			}
		}
		rand.Shuffle(len(msg), func(i, j int) { msg[i], msg[j] = msg[j], msg[i] })

		for i := 0; i < numThreads; i++ {
			part := getPart(msg, i, numThreads)
			threadMessages[i] <- part
		}
		fmt.Println("Done startup")
	}
}
