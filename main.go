package main

import (
	"wsserver/bootstrap/socket"
	"wsserver/bootstrap/websocket"
	"wsserver/db"
	"time"
)

func main() {
	db.InitDB()
	go wsinit.InitUDP()
	go wsinit.InitTCP()
	go socketinit.InitUDP()
	go socketinit.InitTCP()
	for {
		time.Sleep(100 * time.Second)
	}
}
