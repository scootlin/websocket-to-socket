package socketinit

import (
	"log"
	"mcpe_websocket/db"
	"mcpe_websocket/socketagent"
	"net"
	"strings"
)

var clients map[string]*net.TCPConn
var mcpeList map[string]map[string]string

func InitUDP() {
	// Receive data from udp socket connection
	srcAddr := db.SetSrcAddr()
	log.Println("Bind udp socket : ", srcAddr)
	udpAddr, _ := net.ResolveUDPAddr("udp", srcAddr)
	//	udpAddr, _ := net.ResolveUDPAddr("udp", "192.168.11.98:12005")
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Println("error: ", err)
		return
	}
	log.Println("UDP socket waiting for clients")
	go socketagent.HandleUdpConnection(udpConn)

	// Send data from udp socket connection
	dstAddr := db.SetDstAddr()
	// log.Println("Connect to udp socket : ", dstAddr)
	//	ip := net.ParseIP("192.168.11.9")
	laddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	raddr, _ := net.ResolveUDPAddr("udp", dstAddr)
	//raddr := &net.UDPAddr{IP: ip, Port: 12008}
	conn, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("UDP socket connects successfully")
	go socketagent.SetUDPConn(conn)
}

func InitTCP() {
	clients = make(map[string]*net.TCPConn)
	mcpeList = make(map[string]map[string]string)
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:9090")
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Println(err)
		return
	}
	defer tcpListener.Close()

	log.Println("TCP socket waiting for clients")

	for {
		conn, err := tcpListener.AcceptTCP()
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()

		address := conn.RemoteAddr().String()
		var mcpeName, mcpeNum string
		ip := strings.Split(address, ":")[0]
		if ip != "127.0.0.1" {
			mcpeName, mcpeNum = db.SetMcpeInfo(ip)
			clients[mcpeNum] = conn
			mcpeList[mcpeNum] = make(map[string]string)
			mcpeList[mcpeNum]["number"] = mcpeNum
			mcpeList[mcpeNum]["name"] = mcpeName
		} else if ip == "127.0.0.1" {
			mcpeNum = "localhost"
			clients["localhost"] = conn
		}

		log.Printf("tcp %v is connected", mcpeNum)
		go socketagent.HandleTcpConnection(conn, clients, mcpeList, mcpeNum)
	}
}
