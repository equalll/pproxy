package serve
import "github.com/equalll/mydebug"

import (
	"log"
	"time"
)

type clientSession struct {
	Ip               string
	Port             string
	RequestNum       int
	FirstRequestTime time.Time
	LastRequestTime  time.Time
	User             *User
}

func (ser *ProxyServe) regirestReq(reqCtx *requestCtx) {mydebug.INFO()
	ip := reqCtx.GetIp()
	now := time.Now()
	ser.mu.Lock()
	defer ser.mu.Unlock()
	var session *clientSession
	client, has := ser.ProxyClients[ip]
	if has {
		session = client
	} else {
		session = &clientSession{
			Ip:               ip,
			RequestNum:       0,
			FirstRequestTime: now,
			LastRequestTime:  now,
		}
	}
	if reqCtx.User.Name == "" && session.User != nil {
		reqCtx.User = session.User
	} else if reqCtx.User.Name != "" {
		session.User = reqCtx.User
	}

	session.LastRequestTime = now
	session.RequestNum++
	if ser.Debug {
		log.Println("session_debug:", session)
	}
	ser.ProxyClients[ip] = session

	reqCtx.ClientSession = session

	if !has {
		ser.wsSer.broadProxyClientNum()
	}
}

func (ser *ProxyServe) cleanExpiredSession() {mydebug.INFO()
	ser.mu.Lock()
	defer ser.mu.Unlock()
	now := time.Now()
	deleteIps := []string{}
	for ip, session := range ser.ProxyClients {
		t := now.Sub(session.LastRequestTime)
		if t.Minutes() > 10 {
			deleteIps = append(deleteIps, ip)
		}
	}
	for _, ip := range deleteIps {
		delete(ser.ProxyClients, ip)
		log.Println("session expired:ip=", ip)
	}
	ser.wsSer.broadProxyClientNum()
}
