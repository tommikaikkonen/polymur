// Package listener handles the ingestion for polymur
package listener

import (
	"log"
	"net"
	"time"

	"github.com/chrissnell/polymur/statstracker"
)

// UDPListenerConfig holds configuration for the TCP listener
type UDPListenerConfig struct {
	IncomingQueue chan []*string
	FlushTimeout  int
	FlushSize     int
	Stats         *statstracker.Stats
	IP            string
	Port          int
	Zone          string
}

// UDPListener Listens for messages.
func UDPListener(config *UDPListenerConfig) {
	log.Printf("Metrics listener started: %s\n", config.Addr)
	server, err := net.ListenUDP("udp", config.Addr)
	if err != nil {
		log.Fatalf("Listener error: %s\n", err)
	}

	messages := make(chan string, 128)
	go messageBatcher(messages, config)

	buf := make([]byte, 1024)
	go func() {
		for {
			n, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				log.Printf("UDP read error: %s\n", err)
				continue
			}
			// We need to have each message in its own
			// array so that the slices sent to messages don't
			// have their data changed underneath if a new packet
			// is received before processing the previous one.
			msg = make([]string, n)
			copy(msg, buf[:n])
			messages <- msg
		}
		defer server.Close()
	}()
}

func messageBatcher(messages chan string, config *UDPListenerConfig) {
	flushTimeout := time.NewTicker(time.Duration(config.FlushTimeout) * time.Second)
	defer flushTimeout.Stop()

	batch := make([]*string, config.FlushSize)
	pos := 0

run:
	for {
		// We hit the flush timeout, load the current batch if present.
		select {
		case <-flushTimeout.C:
			if len(batch) > 0 {
				config.IncomingQueue <- batch
				batch = make([]*string, config.FlushSize)
				pos = 0
			}
			// reset timeout
			flushTimeout := time.NewTicker(time.Duration(config.FlushTimeout) * time.Second)
		case m, ok := <-messages:
			if !ok {
				break run
			}

			// Drop message and respond if the incoming queue is at capacity.
			if len(config.IncomingQueue) >= 32768 {
				log.Printf("Incoming queue capacity %d reached\n", 32768)
				// Needs some flow control logic.
			}

			// If this puts us at the FlushSize threshold, enqueue
			// into the q.
			if pos+1 >= config.FlushSize {
				batch[config.FlushSize-1] = &m
				config.IncomingQueue <- batch
				batch = make([]*string, config.FlushSize)
				pos = 0
			} else {
				// Otherwise, just append message to current batch.
				batch[pos] = &m
				pos++
			}
		}
	}

	// Load any partial batch before
	// we return.
	config.IncomingQueue <- batch
}
