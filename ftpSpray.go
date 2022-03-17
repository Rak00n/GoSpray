package main

import (
	"fmt"
	"github.com/GoSpray/goftp"
	"sync"
)



func ftpSpray(wg *sync.WaitGroup, channelToCommunicate chan string,  taskToRun task) {
	defer wg.Done()
	if taskToRun.target.port == 0 {
		taskToRun.target.port = 21
	}
	for _,username := range taskToRun.usernames {
		for _,password := range taskToRun.passwords {

			ftpClient, err := goftp.NewFtp(stringifyTarget(taskToRun.target))
			if err != nil {
				panic(err)
			}

			if err = ftpClient.Login(username, password); err != nil {
				fmt.Print("-")
			} else {
				fmt.Print("+")
				channelToCommunicate <- username+":"+password
			}
			ftpClient.Close()
		}
	}


}