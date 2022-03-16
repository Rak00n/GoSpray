package main

import (
	// "os"
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"
	// "FtpSpray"
)

var pathToUsernameList string
var usernameListRandomization bool
var pathToPasswordList string
var passwordListRandomization bool
var protocol string
var target string
var workersNumber int
var taskStateObj taskState 

func init() {
	flag.StringVar(&pathToUsernameList, "ul", "usernames.txt", "Path to usernames list")
	flag.StringVar(&pathToPasswordList, "pl", "passwords.txt", "Path to passwords list")
	flag.BoolVar(&usernameListRandomization, "ru", true, "Randomize users list 1/0")
	flag.BoolVar(&passwordListRandomization, "rp", true, "Randomize passwords list 1/0")
	flag.StringVar(&protocol, "p", "ftp", "Protocol")
	flag.StringVar(&target, "t", "10.0.0.1:21", "Target")
	flag.IntVar(&workersNumber, "w", 5, "Number of Workers")
	flag.Parse()
	
	
}

func printSuccessfulLogin(c chan string) {
	credentials := <- c
	fmt.Println("\nSuccess: "+credentials)
}

func main() {
	randomSeed := time.Now().UnixNano()
	taskStateObj.taskRandomSeed = randomSeed
	usernames := loadList(pathToUsernameList)
	passwords := loadList(pathToPasswordList)
	rand.Seed(taskStateObj.taskRandomSeed)

    targetToSpray := parseTarget(target)
	wholeTask := task{target: targetToSpray, usernames: usernames, passwords: passwords, numberOfWorkers: workersNumber}
	tasks := dispatchTask(wholeTask)

	var wg sync.WaitGroup
	channelForWorker := make(chan string)
	go printSuccessfulLogin(channelForWorker)
	if protocol == "ftp" {
		for _,task := range tasks{
			wg.Add(1)

			go ftpSpray(&wg,channelForWorker,task)
		}
	} else if protocol == "ssh" {
		for _,task := range tasks{
			wg.Add(1)

			go sshSpray(&wg,channelForWorker,task)
		}
	}

	wg.Wait()
	close(channelForWorker)
}