package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	var (
		g   errgroup.Group
		err error
		end = make(chan bool, 1)
	)

	g.Go(keepAlive(end))
	g.Go(read(end))
	if err = g.Wait(); err != nil {
		log.Printf("Error while doing shit %v", err.Error())
	}
}

func keepAlive(end chan bool) func() error {
	count := 0
	return func() error {
		ticker := time.NewTicker(1 * time.Second)
		var err error
	loop:
		for {
			select {
			case <-end:
				log.Println("Got end signal 2")
				break loop
			case <-ticker.C:
				count++
				if count > 3 {
					log.Println("keepAlive error")
					end <- true
					log.Println("Send end to channel")
					ticker.Stop()
					break loop
				}
			}
		}
		return err
	}
}

func read(end chan bool) func() error {
	var (
		count = 0
	)
	return func() error {
		for {
			select {
			case <-end:
				log.Println("Got end signal1")
				return nil
			default:
				if count > 4 {
					log.Println("read error")
					end <- true
					return fmt.Errorf("read error")
				}
			}
		}
	}
}
