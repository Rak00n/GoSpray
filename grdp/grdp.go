package grdp

import (
	"./core"
	"./glog"
	"./protocol/nla"
	"./protocol/pdu"
	"./protocol/rfb"
	"./protocol/sec"
	"./protocol/t125"
	"./protocol/tpkt"
	"./protocol/x224"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

const (
	PROTOCOL_RDP = "PROTOCOL_RDP"
	PROTOCOL_SSL = "PROTOCOL_SSL"
)

type Client struct {
	Host string // ip:port
	tpkt *tpkt.TPKT
	x224 *x224.X224
	mcs  *t125.MCSClient
	sec  *sec.Client
	pdu  *pdu.Client
	vnc  *rfb.RFB
}

func NewClient(host string, logLevel glog.LEVEL) *Client {
	glog.SetLevel(logLevel)
	logger := log.New(os.Stdout, "", 0)
	glog.SetLogger(logger)
	return &Client{
		Host: host,
	}
}

func (g *Client) LoginForSSL(domain, user, pwd string) error {
	conn, err := net.DialTimeout("tcp", g.Host, 3*time.Second)
	if err != nil {
		return fmt.Errorf("[dial err] %v", err)
	}
	defer conn.Close()
	glog.Info(conn.LocalAddr().String())

	g.tpkt = tpkt.New(core.NewSocketLayer(conn), nla.NewNTLMv2(domain, user, pwd))
	g.x224 = x224.New(g.tpkt)
	g.mcs = t125.NewMCSClient(g.x224)
	g.sec = sec.NewClient(g.mcs)
	g.pdu = pdu.NewClient(g.sec)

	g.sec.SetUser(user)
	g.sec.SetPwd(pwd)
	g.sec.SetDomain(domain)

	g.tpkt.SetFastPathListener(g.sec)
	g.sec.SetFastPathListener(g.pdu)
	g.pdu.SetFastPathSender(g.tpkt)

	err = g.x224.Connect()
	if err != nil {
		return fmt.Errorf("[x224 connect err] %v", err)
	}
	glog.Info("wait connect ok")
	wg := &sync.WaitGroup{}
	breakFlag := false
	wg.Add(1)

	g.pdu.On("error", func(e error) {
		err = e
		glog.Error("error", e)
		g.pdu.Emit("done")
	})
	g.pdu.On("close", func() {
		err = errors.New("close")
		glog.Info("on close")
		g.pdu.Emit("done")
	})
	g.pdu.On("success", func() {
		err = nil
		glog.Info("on success")
		g.pdu.Emit("done")
	})
	g.pdu.On("ready", func() {
		glog.Info("on ready")
		g.pdu.Emit("done")
	})
	g.pdu.On("update", func(rectangles []pdu.BitmapData) {
		glog.Info("on update:", rectangles)
	})
	g.pdu.On("done", func() {
		if breakFlag == false {
			breakFlag = true
			wg.Done()
		}
	})
	wg.Wait()
	return err
}

func (g *Client) LoginForRDP(domain, user, pwd string) error {
	conn, err := net.DialTimeout("tcp", g.Host, 3*time.Second)
	if err != nil {
		return fmt.Errorf("[dial err] %v", err)
	}
	defer conn.Close()
	glog.Info(conn.LocalAddr().String())

	g.tpkt = tpkt.New(core.NewSocketLayer(conn), nla.NewNTLMv2(domain, user, pwd))
	g.x224 = x224.New(g.tpkt)
	g.mcs = t125.NewMCSClient(g.x224)
	g.sec = sec.NewClient(g.mcs)
	g.pdu = pdu.NewClient(g.sec)

	g.sec.SetUser(user)
	g.sec.SetPwd(pwd)
	g.sec.SetDomain(domain)

	g.tpkt.SetFastPathListener(g.sec)
	g.sec.SetFastPathListener(g.pdu)
	g.pdu.SetFastPathSender(g.tpkt)

	g.x224.SetRequestedProtocol(x224.PROTOCOL_RDP)

	err = g.x224.Connect()
	if err != nil {
		return fmt.Errorf("[x224 connect err] %v", err)
	}
	glog.Info("wait connect ok")
	wg := &sync.WaitGroup{}
	breakFlag := false
	updateCount := 0
	wg.Add(1)

	g.pdu.On("error", func(e error) {
		err = e
		glog.Error("error", e)
		g.pdu.Emit("done")
	})
	g.pdu.On("close", func() {
		err = errors.New("close")
		glog.Info("on close")
		g.pdu.Emit("done")
	})
	g.pdu.On("success", func() {
		err = nil
		glog.Info("on success")
		g.pdu.Emit("done")
	})
	g.pdu.On("ready", func() {
		glog.Info("on ready")
	})
	g.pdu.On("update", func(rectangles []pdu.BitmapData) {
		glog.Info("on update:", rectangles)
		updateCount += 1
		//fmt.Println(updateCount," ",rectangles[0].BitmapLength)
	})
	g.pdu.On("done", func() {
		if breakFlag == false {
			breakFlag = true
			wg.Done()
		}
	})

	//wait 2 Second
	time.Sleep(time.Second * 3)
	if breakFlag == false {
		breakFlag = true
		wg.Done()
	}
	wg.Wait()

	if updateCount > 50 {
		return nil
	}
	err = errors.New("login failed")
	return err
}

func Login(target, domain, username, password string) error {
	var err error
	g := NewClient(target, glog.NONE)
	//SSL协议登录测试
	err = g.LoginForSSL(domain, username, password)
	if err == nil {
		return nil
	}
	if err.Error() != PROTOCOL_RDP {
		return err
	}
	//RDP协议登录测试
	err = g.LoginForRDP(domain, username, password)
	if err == nil {
		return nil
	} else {
		return err
	}
}

func LoginForSSL(target, domain, username, password string) error {
	var err error
	g := NewClient(target, glog.NONE)
	//SSL协议登录测试
	err = g.LoginForSSL(domain, username, password)
	if err == nil {
		return nil
	}
	return err
}

func LoginForRDP(target, domain, username, password string) error {
	var err error
	g := NewClient(target, glog.NONE)
	//SSL协议登录测试
	err = g.LoginForRDP(domain, username, password)
	if err == nil {
		return nil
	}
	return err
}

func VerifyProtocol(target string) string {
	var err error
	err = LoginForSSL(target, "", "administrator", "test")
	if err == nil {
		return PROTOCOL_SSL
	}
	if err.Error() != PROTOCOL_RDP {
		return PROTOCOL_SSL
	}
	return PROTOCOL_RDP
}

var target string
var domain string
var username string
var password string

func init() {
	flag.StringVar(&target, "target", "", "Target host and port")
	flag.StringVar(&domain, "domain", "", "User Domain")
	flag.StringVar(&username, "username", "", "Username to try")
	flag.StringVar(&password, "password", "", "Password to try")
	flag.Parse()
}

func main() {
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