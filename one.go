package main

import (
	"fmt"
	"time"
)

func main() {
	defer func(t time.Time) { fmt.Printf("total took %v\n", time.Since(t)) }(time.Now())

	t := time.Now()
	payloads := createPayloads("hello %s, your word of the day is %s")
	fmt.Printf("payload creation took %v\n", time.Since(t))

	t = time.Now()
	sendPayloads(payloads)
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

func sendPayloads(payloads map[string][]string) {
	for platform, payloads := range payloads {
		for _, payload := range payloads {
			switch platform {
			case "ios":
				time.Sleep(100 * time.Millisecond) // network overhead
				fmt.Printf("apns | %s\n", payload)
			case "android":
				time.Sleep(100 * time.Millisecond) // network overhead
				fmt.Printf("gcm | %s\n", payload)
			}
		}
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
