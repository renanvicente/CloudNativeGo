package main

import (
	"fmt"
	"time"
)

func timely() {
	timer := time.NewTimer(5 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()			// Be sure to stop the ticker! or you may have a memory leak

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Tick!")
			case <-done:
				return
			}
		}
	}()

	<-timer.C
	fmt.Println("It's time!")
	close(done)
}

func main() {
	timely()
}
