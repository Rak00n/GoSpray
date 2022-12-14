package main

import 	"github.com/jcmturner/gokrb5/client"

func kerberosSpray() {
	cl := client.NewClientWithPassword("username", "REALM.COM", "password")

}