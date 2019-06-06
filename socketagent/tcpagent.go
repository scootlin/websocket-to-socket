package socketagent

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var tcpClients map[string]*net.TCPConn
var tar = make(map[string]string)
var tlb = make(map[string][]byte)
var rlb = make(map[string][]byte)
var alb = make(map[string][]byte)

type Connlist struct {
	Cat      string                       `json:"cat"`
	Type     string                       `json:"type"`
	Connlist map[string]map[string]string `json:"mcpelist"`
}

func HandleTcpConnection(conn *net.TCPConn, connections map[string]*net.TCPConn, list map[string]map[string]string, mcpeNum string) {
	tcpClients = connections
	connlist := &Connlist{}
	connlist.Cat = "data"
	connlist.Type = "mcpelist_data"
	connlist.Connlist = list
	sendTcpws("65535", connlist)
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Printf("tcp %v is disconnected\n", mcpeNum)
				delete(connections, mcpeNum)
				delete(list, mcpeNum)
				connlist.Connlist = list
				sendTcpws("65535", connlist)
				break
			}
		}
		if buffer != nil {
			buf := buffer[:n]
			if len(buf) < 44 {
				log.Println("buffer ll is not enough , size: ", len(buf))
				continue
			}
			v := uint8(buf[0])
			op := uint8(buf[1])
			ct := uint8(buf[2])
			cle := bytes.Trim(buf[3:19], "\x00")
			cle = bytes.Trim(cle, "'")
			clr := bytes.Trim(buf[19:35], "\x00")
			clr = bytes.Trim(clr, "'")
			ses := binary.LittleEndian.Uint32(buf[35:39])
			dt := binary.LittleEndian.Uint16(buf[39:41])
			pt := uint8(buf[41])
			ll := binary.LittleEndian.Uint16(buf[42:44])
			fmt.Printf("R: v = %v op = %v ct = %v cle = %v clr = %v ses = %v dt = %v pt = %v ll = %v\n", v, op, ct, cle, clr, ses, dt, pt, ll)

			pack := Packet{
				vs: v,
				o:  op,
				c:  ct,
				s:  ses,
				d:  dt,
				p:  pt,
				l:  ll,
			}
			pack.Callee = string(cle)
			pack.Caller = string(clr)

			if ll > 0 {
				pack.Payload = buf[44:]
			}
			switch op {
			case 1:
				key := string(cle) + string(clr)
				timestamp := binary.LittleEndian.Uint32(buf[44:48])
				t := time.Unix(int64(timestamp), 0)
				tar[key] = t.Format("2006-01-02 15:04:05")
			case 4:
				key := string(cle) + string(clr)
				t := time.Now().Format("2006_01_02_150405")
				switch dt {
				case 20742:
					filename := "Tl" + t + ".csv"
					saveCSVfile(tlb[key], filename, string(clr), dt, t)
					delete(tlb, key)
				case 20749:
					filename := "Rl" + t + ".csv"
					saveCSVfile(rlb[key], filename, string(clr), dt, t)
					delete(rlb, key)
				case 20752:
					filename := "Al" + t + ".csv"
					saveCSVfile(alb[key], filename, string(clr), dt, t)
					delete(alb, key)
				}
			case 8:
				switch dt {
				case 9888:
					func9888(pack, mcpeNum)
				default:
					doFunc(dt, pack)
				}
			}
		}
	}
}

func sendTcpws(cle string, data interface{}) {
	if Tcpws != nil {
		jsondata, err := json.Marshal(data)
		if err != nil {
			log.Println("error: ", err)
		}
		Tcpws(cle, string(jsondata))
	} else {
		log.Println("there is not any tcp webosocket connected")
	}
}
