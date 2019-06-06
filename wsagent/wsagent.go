package wsagent

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"wsserver/socketagent"

	"github.com/gorilla/websocket"
)

var tcpClients = make(map[string]*websocket.Conn)
var udpClients = make(map[string]*websocket.Conn)
var funcRoute = map[string]func(buffer []byte){}

func HandleUdpConnection(w http.ResponseWriter, r *http.Request) {
	socketagent.SendToUdpws(sendToUdpws)
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}
	sr := r.URL.Query()["sr"][0]
	fmt.Printf("sr = %v\n", sr)
	if _, ok := udpClients[sr]; ok {
		fmt.Printf("ID existed, remove sr = %v\n", sr)
		jsondata, err := json.Marshal(logout)
		if err != nil {
			log.Println("error: ", err)
		}
		sendToUdpws(sr, string(jsondata))
		err = udpClients[sr].Close()
		if err != nil {
			fmt.Println("Close client connect error : %s", err.Error())
		}
		delete(udpClients, sr)
	}
	udpClients[sr] = conn
	addr := strings.Split(r.Header.Get("Origin"), "http://")
	log.Printf("udpws %v is connected", addr[1])
	go func() {
		defer func() {
			if _, ok := udpClients[sr]; ok {
				delete(udpClients, sr)
			}
			conn.Close()
		}()
		for {
			_, buffer, err := conn.ReadMessage()
			if err != nil {
				log.Print(err)
				return
			}
			if buffer != nil {
				jsond := make(map[string]interface{})
				err := json.Unmarshal(buffer, &jsond)
				// log.Print("Get websocket command : ", jsond)
				if err != nil {
					log.Print(err)
					return
				}
				if fn, ok := jsond["type"]; ok {
					doFunc(fn.(string), buffer)
				}
				jsond = nil
			}
		}
	}()
}

func HandleTcpConnection(w http.ResponseWriter, r *http.Request) {
	socketagent.SendToTcpws(sendToTcpws)
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}
	go func() {
		defer conn.Close()
		for {
			_, buffer, err := conn.ReadMessage()
			if err != nil {
				log.Print(err)
				return
			}
			if buffer != nil {
				jsond := make(map[string]interface{})
				err := json.Unmarshal(buffer, &jsond)
				log.Print("Get websocket command : ", jsond)
				if err != nil {
					log.Print(err)
					return
				}
				if fn, ok := jsond["type"]; ok {
					doFunc(fn.(string), buffer)
				}
				jsond = nil
			}
		}
	}()
}

func doFunc(t string, buffer []byte) {
	if fn, ok := funcRoute[t]; ok {
		fn(buffer)
	}
}

func sendToUdpws(index string, buf string) {
	if index == "65535" {
		for _, c := range udpClients {
			websocket.WriteJSON(c, buf)
		}
	} else {
		if conn, ok := udpClients[index]; ok {
			websocket.WriteJSON(conn, buf)
		} else {
			log.Printf("udpws %v is disconnected\n", index)
		}
	}
}

func sendToTcpws(index string, buf string) {
	if index == "65535" {
		for _, c := range tcpClients {
			websocket.WriteJSON(c, buf)
		}
	} else {
		if conn, ok := tcpClients[index]; ok {
			websocket.WriteJSON(conn, buf)
		} else {
			log.Printf("tcpws %v is disconnected\n", index)
		}
	}
}
