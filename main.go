package main

import (
	"bufio"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	VERSION = "1.0"
)

var (
	CfAddr        = ""
	AllowHostMap  = []string{}
	ListenPort    = 80
	ListenAddress = ""
	Debug         = false
)

func init() {
	var allowHost string
	fmt.Printf("cfp version %s\n", VERSION)

	flag.StringVar(&CfAddr, "cfaddr", "", "cloudflare node address")
	flag.StringVar(&allowHost, "hosts", "", "allow host(s),separate by comma")
	flag.StringVar(&ListenAddress, "addr", "0.0.0.0", "listen addr")
	flag.IntVar(&ListenPort, "port", 80, "listen port")
	flag.BoolVar(&Debug, "debug", false, "show debug logs")

	flag.Parse()

	if CfAddr == "" || allowHost == "" || ListenAddress == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	fmt.Printf("cfaddr: %s\n", CfAddr)
	fmt.Printf("hosts: %s\n", allowHost)
	fmt.Printf("addr: %s\n", ListenAddress)
	fmt.Printf("port: %d\n", ListenPort)
	fmt.Printf("debug: %v\n", Debug)

	if Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	hosts := strings.Split(allowHost, ",")
	for _, v := range hosts {
		AllowHostMap = append(AllowHostMap, v)
	}

}
func main() {
	listenAddr := fmt.Sprintf("%s:%d", ListenAddress, ListenPort)
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("listen at %s failed: %v\n", listenAddr, err)
	}

	log.Println("server running at", listenAddr)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Errorf("accept client failed: %v\n", err)
			log.Errorln("sleep one second and retry")
			time.Sleep(time.Second * 1)
			continue
		}
		go handle(conn)
	}
}

func checkHost(t string) bool {
	t = strings.Split(t, ":")[0]
	for _, v := range AllowHostMap {
		if strings.HasSuffix(t, v) {
			return true
		}
	}
	return false
}

func handle(conn net.Conn) {
	var ok bool
	defer func() {
		if !ok {
			conn.Close()
		}
	}()

	log.Debugf("new client coming: %s", conn.RemoteAddr())

	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		log.Debugf("read client http request failed: %s, %v", conn.RemoteAddr(), err)
		return
	}

	if req.Header.Get("cfp") != "" {
		log.Errorf("found loop: %s, %v", conn.RemoteAddr(), err)
		return
	}

	if !checkHost(req.Host) {
		log.Debugf("client host mismatch: %s, %s", conn.RemoteAddr(), req.Host)
		return
	}

	addr := CfAddr
	if addr == "direct" {
		req.Header.Add("cfp", "1")
		addr = req.Host
		if !strings.Contains(addr, ":") {
			addr += ":80"
		}
	} else if len(strings.Split(CfAddr, ":")) != 2 {
		port := "80"
		portArr := strings.Split(req.Host, ":")
		if len(portArr) == 2 {
			port = portArr[1]
		}
		addr = addr + ":" + port
	}

	proxy, err := net.Dial("tcp", addr)
	if err != nil {
		log.Errorf("dial cloudflare node failed: %v", err)
		return
	}

	log.Debugf("proxy connection established : %s, %s", conn.RemoteAddr(), req.Host)

	ok = true

	req.Write(proxy)

	defer proxy.Close()

	go io.Copy(conn, proxy)
	io.Copy(proxy, conn)

}
