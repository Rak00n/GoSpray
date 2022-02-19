package main

import (
	"fmt"
	// "os"
	"flag"
	"time"
	"math/rand"
	// "FtpSpray"
)

var pathToUsernameList string
var usernameListRandomization bool
var pathToPasswordList string
var passwordListRandomization bool
var protocol string
var target string
var workersNumber int

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

func main() {
	rand.Seed(time.Now().UnixNano())
    fmt.Println("hello world")
    fmt.Println(pathToUsernameList)
    fmt.Println(pathToPasswordList)
    fmt.Println(protocol)
    fmt.Println(target)
    targetToSpray := parseTarget(target)
	if protocol == "ftp" {
		ftpSpray(targetToSpray)
	}
	
}