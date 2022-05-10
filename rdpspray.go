package main

import (
	"github.com/Rak00n/grdp"
	"log"
	"os"
)

func rdpSpray () {
	glog.SetLevel(glog.INFO)
	logger := log.New(os.Stdout, "", 0)
	glog.SetLogger(logger)
	g := rdp.NewClient(ip, glog.LEVEL(loglevel))
	err := g.Login(domain, user, passwd)
}
