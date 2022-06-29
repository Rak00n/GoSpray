package main

import (
	"fmt"

	"github.com/GoSpray/grdp/glog"
	"github.com/GoSpray/grdp/protocol/pdu"
	"github.com/GoSpray/grdp/protocol/sec"
	"github.com/GoSpray/grdp/protocol/t125"
	"github.com/GoSpray/grdp/protocol/tpkt"
	"github.com/GoSpray/grdp/protocol/x224"
)

type Client struct {
	Host string // ip:port
	tpkt *tpkt.TPKT
	x224 *x224.X224
	mcs  *t125.MCSClient
	sec  *sec.Client
	pdu  *pdu.Client
}



func rdpSpray () {
	client := NewClient(target, glog.NONE)
	var err error
	//if useNLA {
	//	err = client.LoginForSSL(domain,username, password)
	//} else {
	//	err = client.LoginForRDP(domain,username, password)
	//}
	err = client.LoginForSSL(domain,username, password)
	if err != nil {
		fmt.Println("login failed,", err)
	} else {
		fmt.Println("login success")
	}
}
