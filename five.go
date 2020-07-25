package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func collectLogs(wg *sync.WaitGroup, closeChan chan struct{}) (chan<- bool, <-chan map[bool]int) {
	logChan := make(chan bool, 100)
	resultChan := make(chan map[bool]int)
	go func() {
		log := make(map[bool]int)
		for {
			select {
			case success := <-logChan:
				log[success]++
				wg.Done()
			case <-closeChan:
				resultChan <- log
				return
			}
		}
	}()
	return logChan, resultChan
}

func sendMessages(logChan chan<- bool) chan<- string {
	msgChan := make(chan string) // limiter will be the blocker anyway
	go func() {
		limiter := make(chan struct{}, 1000)
		for message := range msgChan {
			_ = message
			limiter <- struct{}{}
			go func() {
				time.Sleep(100 * time.Millisecond)
				if rand.Intn(10) == 0 {
					logChan <- false
				} else {
					logChan <- true
				}
				<-limiter
			}()
		}
	}()
	return msgChan
}

func main() {
	defer func(t time.Time) { fmt.Printf("total took %v\n", time.Since(t)) }(time.Now())

	wg := new(sync.WaitGroup)

	closeChan := make(chan struct{})

	logChan, resultChan := collectLogs(wg, closeChan)

	msgChan := sendMessages(logChan)

	for _, message := range generateDummyMessages() {
		wg.Add(1)
		msgChan <- message
	}

	wg.Wait()

	close(closeChan)

	fmt.Println(<-resultChan)
}

func generateDummyMessages() []string {
	messages := make([]string, 10000)
	for i := range messages {
		messages[i] = strconv.Itoa(i)
	}
	return messages
}
