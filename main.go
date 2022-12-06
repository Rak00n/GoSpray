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
	flag.BoolVar(&usernameListRandomization, "ru", false, "Randomize users list")
	flag.BoolVar(&passwordListRandomization, "rp", false, "Randomize passwords list")
	flag.StringVar(&protocol, "p", "ftp", "Protocol (ftp,ssh,httpbasic,httpdigest,rdp,winldap)")
	flag.StringVar(&target, "t", "10.0.0.1:21", "Target")
	flag.IntVar(&workersNumber, "w", 5, "Number of Workers")
	flag.Parse()
	
	
}

func printSuccessfulLogin(c chan string) {
	for {
		credentials := <- c
		fmt.Println("\nSuccess: "+credentials)
	}

}

type runningTask struct {
	RandomSeed int64
	UsersList string
	PasswordsList string
	ProtocolToSpray string
	Target string
	WorkersCount int
	WorkersStates []workerState
	UsernamesRandomization bool
	PasswordsRandomization bool
}

var currentTask runningTask

func main() {
	if restoreTask == true {
		err := readGob("./progress.gob",&currentTask)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		if usernameListRandomization || passwordListRandomization {
			currentTime := time.Now().UnixNano()
			currentTask.RandomSeed = currentTime
			if usernameListRandomization {
				currentTask.UsernamesRandomization = true
			}
			if passwordListRandomization {
				currentTask.PasswordsRandomization = true
			}
		} else {
			currentTask.RandomSeed = 0
			currentTask.PasswordsRandomization = false
			currentTask.UsernamesRandomization = false
		}


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


	usernames := loadList(currentTask.UsersList)
	passwords := loadList(currentTask.PasswordsList)

	if currentTask.UsernamesRandomization{
		rand.Seed(currentTask.RandomSeed)
		rand.Shuffle(len(usernames), func(i, j int) { usernames[i], usernames[j] = usernames[j], usernames[i] })
	}
	if currentTask.PasswordsRandomization{
		rand.Seed(currentTask.RandomSeed)
		rand.Shuffle(len(passwords), func(i, j int) { passwords[i], passwords[j] = passwords[j], passwords[i] })
	}

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
	} else if currentTask.ProtocolToSpray == "rdp" {

		for _,task := range tasks{
			wg.Add(1)
			go rdpSpray(&wg,channelForWorker,task,&currentTask.WorkersStates[iter].WorkerProgress)
			iter++
		}
	} else if currentTask.ProtocolToSpray == "winldap" {

		for _,task := range tasks{
			wg.Add(1)
			go ldapSpray(&wg,channelForWorker,task,&currentTask.WorkersStates[iter].WorkerProgress)
			iter++
		}
	}
	go monitorCurrentTask()
	wg.Wait()
	close(channelForWorker)
}


