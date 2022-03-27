package main

import (
	"github.com/GoSpray/rdp"
)

func rdpSpray () {
	g := rdp.NewClient(ip, glog.LEVEL(loglevel))
	err := g.Login(domain, user, passwd)
}
