package main

import (
	"github.com/GoSpray/grdp"
)

func rdpSpray () {
	g := rdp.NewClient(ip, glog.LEVEL(loglevel))
	err := g.Login(domain, user, passwd)
}
