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
	scheme string
	url string
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
	tempString := targetString
	// https://www.www.www:52/sasd/1/2/3
	if strings.Contains(targetString, "://") {
		s := strings.Split(targetString, "://")
		target.scheme = s[0]
		tempString = s[1]
	}
	// www.www.www:52/sasd/1/2/3
	if strings.Contains(tempString, ":") {
		s := strings.Split(tempString, ":")
		target.host = s[0]
		tempString = s[1]
		if strings.Contains(tempString,"/"){
			tempStringSlice := strings.Split(tempString,"/")
			target.port,_ = strconv.Atoi(tempStringSlice[0])
			target.url = strings.Join(tempStringSlice[1:],"/")
		} else {
			target.port,_ = strconv.Atoi(tempString)
		}
	} else {
		if strings.Contains(tempString,"/"){
			tempStringSlice := strings.Split(tempString,"/")
			target.host = tempStringSlice[0]
			target.url = strings.Join(tempStringSlice[1:],"/")
			target.port = 0
		} else {
			target.host = tempString
			target.port = 0
		}
	}

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
