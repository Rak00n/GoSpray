package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz"+"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}



type targetStruct struct{
	host string
	port int
	scheme string
	url string
}

type workerState struct {
	WorkerId int
	WorkerProgress int
}

type taskState struct {
	TaskRandomSeed int64
	WorkersStates []workerState
}

type task struct {
	target targetStruct
	usernames []string
	passwords []string
	numberOfWorkers int
}

func writeGob(filePath string,object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func readGob(filePath string,object interface{}) error {
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
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

func logProgress(logFilename string, logMessage string) {
	f, err := os.OpenFile(logFilename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(logMessage+"\n"); err != nil {
		panic(err)
	}
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

func saveProgress () {
	err := writeGob("./progress.gob",currentTask)
	if err != nil{
		fmt.Println(err)
	}
}

func monitorCurrentTask() {
	for {
		//fmt.Println(currentTask)
		saveProgress()
		time.Sleep(time.Second)
	}
}