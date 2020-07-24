package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	defer func(t time.Time) { fmt.Printf("total took %v\n", time.Since(t)) }(time.Now())

	t := time.Now()
	payloads := createPayloads("hello %s, your word of the day is %s")
	fmt.Printf("payload creation took %v\n", time.Since(t))

	wg := new(sync.WaitGroup)
	wg.Add(len(payloads["ios"]) + len(payloads["android"]))

	apns := newAppleSender(wg)
	gcm := newGoogleSender(wg)

	t = time.Now()
	sendPayloads(apns, gcm, payloads)
	wg.Wait()
	fmt.Printf("payload sending took %v\n", time.Since(t))
}

func createPayloads(message string) map[string][]string {
	payloads := make(map[string][]string)
	iosUsers := getIOSUsers()
	androidUsers := getAndroidUsers()
	words := getWords()
	for i, word := range words {
		payloads["ios"] = append(payloads["ios"], fmt.Sprintf(message, iosUsers[i], word))
		payloads["android"] = append(payloads["android"], fmt.Sprintf(message, androidUsers[i], word))
	}
	return payloads
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

func sendPayloads(apns chan<- string, gcm chan<- string, payloads map[string][]string) {
	for platform, payloads := range payloads {
		for _, payload := range payloads {
			switch platform {
			case "ios":
				apns <- payload
			case "android":
				gcm <- payload
			}
		}
	}
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
