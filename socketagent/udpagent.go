package socketagent

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func HandleUdpConnection(udpConn *net.UDPConn) {
	defer udpConn.Close()
	for {
		buffer := make([]byte, 1024)
		n, _, err := udpConn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("error: ", err)
			break
		}
		if buffer != nil {
			buf := buffer[:n]
			if len(buf) < 44 {
				log.Println("buffer le is not enough , size: ", len(buf))
				continue
			}
			av := uint8(buf[0])
			bb := uint8(buf[1])
			ty := uint8(buf[2])
			cae := binary.LittleEndian.Uint32(buf[3:19])
			bae := binary.LittleEndian.Uint32(buf[19:35])
			io := binary.LittleEndian.Uint32(buf[35:39])
			dat := binary.LittleEndian.Uint16(buf[39:41])
			pr := uint8(buf[41])
			le := binary.LittleEndian.Uint16(buf[42:44])
			fmt.Printf("R: av = %v bb = %v ty = %v cae = %v bae = %v io = %v dat = %v pr = %v le = %v\n", av, bb, ty, cae, bae, io, dat, pr, le)
		}
	}
}

func sendUdpws(cae string, data interface{}) {
	if Udpws != nil {
		jsondata, err := json.Marshal(data)
		if err != nil {
			log.Println("error: ", err)
		}
		Udpws(cae, string(jsondata))
	} else {
		log.Println("there is not any udp webosocket connected")
	}
}
