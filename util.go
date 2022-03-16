package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type targetStruct struct{
	host string
	port int
}

type workerState struct {
	workerId int
	workerProgress int
}

type taskState struct {
	taskRandomSeed int64
	workersStates []workerState
}

type task struct {
	target targetStruct
	usernames []string
	passwords []string
	numberOfWorkers int
}

func parseTarget(targetString string) targetStruct {
	var target targetStruct
	s := strings.Split(targetString, ":")
	target.host = s[0]
	target.port,_ = strconv.Atoi(s[1])
	return target
}

func stringifyTarget(targetToStringify targetStruct) string {
	portString := strconv.Itoa(targetToStringify.port)
	return targetToStringify.host+":"+portString
}

func logProgress() {
	fmt.Println(taskStateObj)
}

func loadList(pathToFile string) []string {
	var loadedItems []string
	file, err := os.Open(pathToFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		loadedItems = append(loadedItems,scanner.Text())
	}
	return loadedItems
}

func dispatchTask(taskToDispatch task) []task{
	var tasksToReturn []task
	totalUsernames := len(taskToDispatch.usernames)
	usernamesPerWorker := int(math.Ceil(float64(totalUsernames)/float64(taskToDispatch.numberOfWorkers)))

	var usernamesBuffer []string
	iterator := 0
	for _, username := range taskToDispatch.usernames{
		usernamesBuffer = append(usernamesBuffer,username)
		iterator++
		if iterator == usernamesPerWorker {
			tasksToReturn = append(tasksToReturn,task{target: taskToDispatch.target, usernames: usernamesBuffer, passwords: taskToDispatch.passwords, numberOfWorkers: 1})
			usernamesBuffer = nil
			iterator = 0
		}
	}
	if len(usernamesBuffer) != 0 {
		tasksToReturn = append(tasksToReturn,task{target: taskToDispatch.target, usernames: usernamesBuffer, passwords: taskToDispatch.passwords, numberOfWorkers: 1})
	}
	return tasksToReturn
}
