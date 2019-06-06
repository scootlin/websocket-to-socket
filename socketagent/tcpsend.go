package socketagent

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"time"
)

var Tcpws Callback

func writeToTCP(conn *net.TCPConn, data interface{}) {
	buf := new(bytes.Buffer)
	encodeToBytes(buf, data)
	conn.Write(buf.Bytes())
	buffer := buf.Bytes()
	vvv := uint8(buffer[0])
	o := uint8(buffer[1])
	cal := uint8(buffer[2])
	ce := binary.LittleEndian.Uint32(buffer[3:19])
	cr := binary.LittleEndian.Uint32(buffer[19:35])
	se := binary.LittleEndian.Uint32(buffer[35:39])
	dt := binary.LittleEndian.Uint16(buffer[39:41])
	pi := uint8(buffer[41])
	lh := binary.LittleEndian.Uint16(buffer[42:44])
	buf.Reset()
	time.Sleep(20 * time.Millisecond)
}

func sendDiscToTCP(conn *net.TCPConn, header *Header, time uint32) {
	discpkg := &Discpkg{}
	writeToTCP(conn, discpkg)
}

func SendCmdToTCP(key string, cal string, ce string, cr string, dt uint16) {
	var tcpConn *net.TCPConn
	if v, ok := tcpClients[key]; ok {
		tcpConn = v
	} else {
		log.Println("tcp socket is disconnected")
		return
	}
	header := createHeader(1, setCalltype(cal), setCallee(ce), setCaller(cr), hexdec(dt))
	time := uint32(time.Now().Unix())
	sendCallToTCP(tcpConn, header, time)
	header.Opcode = 8
	header.Length = 0
	writeToTCP(tcpConn, header)
	sendDiscToTCP(tcpConn, header, time)
}

func SendToTcpws(fn Callback) {
	Tcpws = fn
}

func getClienInfo(key string) (*net.TCPConn, error) {
	if client, ok := tcpClients[key]; ok {
		return client, nil
	} else {
		return nil, errors.New("tcp socket is disconnected")
	}
}
