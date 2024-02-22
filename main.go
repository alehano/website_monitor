package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	checkURL          = os.Getenv("CHECK_URL")
	checkInterval     = time.Duration(os.Getenv("CHECK_INTERVAL")) * time.Minute
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
	body, err := ioutil.ReadAll(resp.Body)
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
	body, _ := ioutil.ReadAll(response.Body)
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
		log.Println("The specified text is still on the page. No action taken.")
	}
}

func main() {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	checkPageForText() // initial check

	for range ticker.C {
		checkPageForText()
	}
}
