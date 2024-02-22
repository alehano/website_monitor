package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	checkURL          = os.Getenv("CHECK_URL")
	checkInterval     = 15 * time.Minute
	telegramAPIKeyVar = os.Getenv("TELEGRAM_API_KEY")
	chatID            = os.Getenv("CHAT_ID")
	textToCheck       = os.Getenv("TEXT_TO_CHECK")
	shouldContain     = os.Getenv("SHOULD_CONTAIN") == "true"
	alertMessage      = os.Getenv("ALERT_MESSAGE")
)

func fetchPageContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func sendMessage(message string) error {
	baseURL := "https://api.telegram.org/bot" + telegramAPIKeyVar + "/sendMessage"
	response, err := http.PostForm(baseURL, url.Values{"chat_id": {chatID}, "text": {message}})
	if err != nil {
		return err
	}
	defer response.Body.Close()
	// Optionally, you can read and log the response from Telegram
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	fmt.Println("Response from Telegram: ", string(body))
	return nil
}

func checkPageForText() {
	content, err := fetchPageContent(checkURL)
	if err != nil {
		log.Printf("Error fetching page content: %v", err)
		return
	}
	if (shouldContain && strings.Contains(content, textToCheck)) || (!shouldContain && !strings.Contains(content, textToCheck)) {
		err := sendMessage(alertMessage)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	} else {
		log.Println("Condition not met")
	}
}

func main() {
	if os.Getenv("CHECK_INTERVAL") != "" {
		checkIntervalD, err := time.ParseDuration(os.Getenv("CHECK_INTERVAL"))
		if err != nil {
			log.Printf("Error parsing CHECK_INTERVAL: %v", err)
		} else {
			checkInterval = checkIntervalD
		}
	}

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	checkPageForText() // initial check

	for range ticker.C {
		checkPageForText()
	}
}
