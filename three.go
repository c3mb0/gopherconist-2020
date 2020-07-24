package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	defer func(t time.Time) { fmt.Printf("total took %v\n", time.Since(t)) }(time.Now())

	wg := new(sync.WaitGroup)

	apns := newAppleSender(wg)
	gcm := newGoogleSender(wg)

	pc := newPayloadCreator(apns, gcm, wg)

	pc.createAndSend("hello %s, your word of the day is %s")

	wg.Wait()
}

func getAndroidUsers() []string {
	time.Sleep(500 * time.Millisecond) // database overhead
	return androidUsers
}

func getIOSUsers() []string {
	time.Sleep(500 * time.Millisecond) // database overhead
	return iosUsers
}

func getWords() []string {
	time.Sleep(500 * time.Millisecond) // database overhead
	return words
}

type appleSender struct {
	in chan string
	wg *sync.WaitGroup
}

func newAppleSender(wg *sync.WaitGroup) chan<- string {
	as := &appleSender{
		in: make(chan string, 8),
		wg: wg,
	}
	go as.sendApples()
	return as.in
}

func (as *appleSender) sendApples() {
	for payload := range as.in {
		time.Sleep(100 * time.Millisecond) // network overhead
		fmt.Printf("apns | %s\n", payload)
		as.wg.Done()
	}
}

type googleSender struct {
	in chan string
	wg *sync.WaitGroup
}

func newGoogleSender(wg *sync.WaitGroup) chan<- string {
	gs := &googleSender{
		in: make(chan string, 8),
		wg: wg,
	}
	go gs.sendGoogles()
	return gs.in
}

func (gs *googleSender) sendGoogles() {
	for payload := range gs.in {
		time.Sleep(100 * time.Millisecond) // network overhead
		fmt.Printf("gcm | %s\n", payload)
		gs.wg.Done()
	}
}

type payloadCreator struct {
	apns chan<- string
	gcm  chan<- string
	wg   *sync.WaitGroup
}

func newPayloadCreator(apns chan<- string, gcm chan<- string, wg *sync.WaitGroup) *payloadCreator {
	return &payloadCreator{
		apns: apns,
		gcm:  gcm,
		wg:   wg,
	}
}

func (pc *payloadCreator) createAndSendIOS(message string, words []string) {
	iosUsers := getIOSUsers()
	for i, word := range words {
		pc.wg.Add(1)
		pc.apns <- fmt.Sprintf(message, iosUsers[i], word)
	}
	pc.wg.Done()
}

func (pc *payloadCreator) createAndSendAndroid(message string, words []string) {
	androidUsers := getAndroidUsers()
	for i, word := range words {
		pc.wg.Add(1)
		pc.gcm <- fmt.Sprintf(message, androidUsers[i], word)
	}
	pc.wg.Done()
}

func (pc *payloadCreator) createAndSend(message string) {
	words := getWords()
	pc.wg.Add(2)
	go pc.createAndSendAndroid(message, words)
	go pc.createAndSendIOS(message, words)
}

var words = []string{
	"lorem",
	"ipsum",
	"dolor",
	"sit",
	"amet",
	"consectetur",
	"adipiscing",
	"elit",
}

var androidUsers = []string{
	"Shawn",
	"Cari",
	"Alara",
	"Poppy-Mae",
	"Lennie",
	"Jack",
	"Amrita",
	"Mylo",
}

var iosUsers = []string{
	"Brooklyn",
	"Fahad",
	"Chardonnay",
	"Callam",
	"Rebekah",
	"Byron",
	"Rikki",
	"Keyan",
}
