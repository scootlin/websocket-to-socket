package wsinit

import (
	"mcpe_websocket/wsagent"
	"net/http"
)

func InitUDP() {
	http.ListenAndServe(":8080", http.HandlerFunc(wsagent.HandleUdpConnection))
}

func InitTCP() {
	http.ListenAndServe(":8081", http.HandlerFunc(wsagent.HandleTcpConnection))
}
