package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
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
	cookies           = os.Getenv("COOKIES") // "cookie1=value1; cookie2=value2"
	username          = os.Getenv("BASIC_AUTH_USERNAME")
	password          = os.Getenv("BASIC_AUTH_PASSWORD")
)

func fetchPageContent(url string) (string, error) {
	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Add headers to the request
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("DNT", "1")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")

	// Add cookies to the request
	setCookies(req, cookies)

	// Add Basic Authentication to the request
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	// Create a custom HTTP client with TLS configuration
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
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
	body, err := ioutil.ReadAll(response.Body)
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

func setCookies(req *http.Request, cookieString string) {
	if cookieString == "" {
		return
	}
	// Add cookies to the request
	cookies := strings.Split(cookieString, ";")
	for _, cookie := range cookies {
		cookieParts := strings.SplitN(strings.TrimSpace(cookie), "=", 2)
		if len(cookieParts) == 2 {
			req.AddCookie(&http.Cookie{Name: cookieParts[0], Value: cookieParts[1]})
		}
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
		// Generate a random delay of 5-10% of the checkInterval
		delay := time.Duration(rand.Int63n(int64(checkInterval)/10-int64(checkInterval)/20) + int64(checkInterval)/20)
		time.Sleep(delay)

		checkPageForText()
	}
}
