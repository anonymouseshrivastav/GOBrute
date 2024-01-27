package main

import (
	"flag"
	"fmt"
	"gobrute/utils"
	"strconv"
)

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
		utils.Usage("Enter username/email")
	}
	if passlist == "" {
		utils.Usage("Enter password list")
	}
	if url == "" {
		utils.Usage("Enter url")
	}

	fmt.Println(utils.WhiteColor, "Attack Info\n")

	utils.PrintStats("Target: ", url)
	utils.PrintStats("Username: ", username)
	utils.PrintStats("Threads: ", strconv.Itoa(threads)+"\n")

	utils.HandleMain(url, username, passlist, threads)
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
