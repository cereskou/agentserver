package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"ditto.co.jp/agentserver/config"
	"ditto.co.jp/agentserver/logger"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
)

//EventMessage -
type EventMessage struct {
	Mode      string    `json:"mode"`
	IP        string    `json:"ip"`
	Port      uint      `json:"port"`
	TimeStamp time.Time `json:"timestam"`
	ID        int64     `json:"jobid"`
	Status    int       `json:"status"`
}

// NewListen creates a new UDP multicast connection on which to broadcast
func NewListen(port uint) (*net.UDPConn, error) {
	address := fmt.Sprintf(":%v", port)
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// NewBroadcast creates a new UDP multicast connection on which to broadcast
func NewBroadcast(port uint) (*net.UDPConn, error) {
	address := fmt.Sprintf("255.255.255.255:%v", port)
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	local, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", local, addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

//RunBroadcast -
func RunBroadcast(mode string, conf *config.Config) {
	conn, err := NewBroadcast(conf.BroadcastPort)
	if err != nil {
		return
	}
	defer conn.Close()

	ticker := time.Tick(5 * time.Second)
	for now := range ticker {
		// t := now.Format("15:04:05.000")
		msg := EventMessage{
			Mode:      mode,
			IP:        conf.LocalHost,
			Port:      conf.Port,
			TimeStamp: now,
		}

		cmdtxt, _ := json.Marshal(&msg)
		conn.Write([]byte(cmdtxt))
	}
}

//RunListener -
func RunListener(conf *config.Config, agents *AgentMap) {
	conn, err := NewListen(conf.BroadcastPort)
	if err != nil {
		return
	}
	defer conn.Close()

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var v EventMessage
	buf := make([]byte, 512)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		logger.Tracef("Received %v from %v\n", string(buf[0:n]), addr)
		if err != nil {
			logger.Error(err)
		}

		conn.WriteToUDP([]byte("ack"), addr)

		//
		reader := bytes.NewReader(buf[0:n])
		err = json.NewDecoder(reader).Decode(&v)
		//Agent add
		if err == nil && v.Mode == "agent" {
			//check
			key := fmt.Sprintf("%v:%v", v.IP, v.Port)
			agents.Set(key, &v)
		}
	}
}

//HealthCheck -
func HealthCheck(conf *config.Config, agents *AgentMap) {
	ticker := time.Tick(5 * time.Second)
	for now := range ticker {
		keys := agents.Keys()
		if len(keys) == 0 {
			continue
		}
		logger.Tracef("Check alive at %v", now)

		wg := new(sync.WaitGroup)
		dels := make([]string, 0)
		for _, k := range keys {
			wg.Add(1)

			go func(key string) {
				defer wg.Done()

				url := fmt.Sprintf("http://%v/health", key)
				logger.Tracef("Check %v ...", key)

				//Alive check
				client := resty.New()
				client.SetHeaders(map[string]string{
					"User-Agent": "s3transfer distributed",
				})
				client.SetTimeout(2 * time.Second)

				resp, _ := client.R().Get(url)

				alive := string(resp.Body())
				if len(alive) > 0 && alive == "agent.alive" {
					logger.Tracef("%v is alive", key)
					return
				}
				logger.Tracef("%v is down", key)
				//削除対象
				dels = append(dels, key)
			}(k)
		}
		wg.Wait()

		if len(dels) > 0 {
			agents.DeleteList(dels)
		}
	}
}
