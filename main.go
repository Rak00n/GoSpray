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

var restoreTask bool
var pathToUsernameList string
var usernameListRandomization bool
var pathToPasswordList string
var passwordListRandomization bool
var protocol string
var target string
var workersNumber int
var taskStateObj taskState 

func init() {
	flag.BoolVar(&restoreTask, "restore", false, "Restore task")
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

type runningTask struct {
	RandomSeed int
	UsersList string
	PasswordsList string
	ProtocolToSpray string
	Target string
	WorkersCount int
	WorkersStates []workerState
}

var currentTask runningTask

func main() {

	if restoreTask == true {
		err := readGob("./progress.gob",&currentTask)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		rand.Seed(time.Now().UnixNano())
		currentTask.RandomSeed = rand.Intn(100)
		currentTask.UsersList = pathToUsernameList
		currentTask.PasswordsList = pathToPasswordList
		currentTask.ProtocolToSpray = protocol
		currentTask.Target = target
		currentTask.WorkersCount = workersNumber

		for i := 1; i <= workersNumber; i++ {
			currentTask.WorkersStates = append(currentTask.WorkersStates,workerState{
				WorkerId: i,
				WorkerProgress: 0,
			})
		}

		saveProgress()
	}
	fmt.Println(currentTask)

	////////////////////////////

	usernames := loadList(currentTask.UsersList)
	passwords := loadList(currentTask.PasswordsList)

    targetToSpray := parseTarget(currentTask.Target)
	wholeTask := task{target: targetToSpray, usernames: usernames, passwords: passwords, numberOfWorkers: currentTask.WorkersCount}
	tasks := dispatchTask(wholeTask)

	var wg sync.WaitGroup
	channelForWorker := make(chan string)
	go printSuccessfulLogin(channelForWorker)
	iter := 0
	if currentTask.ProtocolToSpray == "ftp" {
		for _,task := range tasks{
			wg.Add(1)
			go ftpSpray(&wg,channelForWorker,task,&currentTask.WorkersStates[iter].WorkerProgress)
			iter++
		}
	} else if currentTask.ProtocolToSpray == "ssh" {
		for _,task := range tasks{
			wg.Add(1)
			go sshSpray(&wg,channelForWorker,task,&currentTask.WorkersStates[iter].WorkerProgress)
			iter++
		}
	} else if currentTask.ProtocolToSpray == "httpbasic" {
		for _,task := range tasks{
			wg.Add(1)
			go basicSpray(&wg,channelForWorker,task,&currentTask.WorkersStates[iter].WorkerProgress)
			iter++
		}
	} else if currentTask.ProtocolToSpray == "httpdigest" {

		for _,task := range tasks{
			wg.Add(1)
			go digestSpray(&wg,channelForWorker,task,&currentTask.WorkersStates[iter].WorkerProgress)
			iter++
		}
	}
	go monitorCurrentTask()
	wg.Wait()
	close(channelForWorker)
}