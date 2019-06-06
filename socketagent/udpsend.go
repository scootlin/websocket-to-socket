package socketagent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

func writeToUDP(data interface{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("writeToUDP error : ", err)
			os.Exit(3)
		}
	}()
	buf := new(bytes.Buffer)
	encodeToBytes(buf, data)
	UdpConn.Write(buf.Bytes())
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
	buf.Reset()
	time.Sleep(200 * time.Microsecond)
}

func SetUDPConn(conn *net.UDPConn) {
	UdpConn = conn
}

func SendToUdpws(fn Callback) {
	Udpws = fn
}

func createVoiceHeader(length uint32, channels uint16, bits uint16, rate uint16) []byte {
	buf := new(bytes.Buffer)
	header := []byte("RIFF")
	binary.Write(buf, binary.LittleEndian, length+36)
	header = append(header, buf.Bytes()...)
	buf.Reset()
	header = append(header, []byte("WAVE")...)
	header = append(header, []byte("fmt ")...)

	buffer := make([]byte, 20)
	binary.LittleEndian.PutUint32(buffer[0:4], uint32(16))
	binary.LittleEndian.PutUint16(buffer[4:6], uint16(1))
	binary.LittleEndian.PutUint16(buffer[6:8], channels)
	binary.LittleEndian.PutUint32(buffer[8:12], uint32(rate))
	binary.LittleEndian.PutUint32(buffer[12:16], uint32(rate*channels*(bits>>3)))
	binary.LittleEndian.PutUint16(buffer[16:18], channels*(bits>>3))
	binary.LittleEndian.PutUint16(buffer[18:20], bits)
	binary.Write(buf, binary.LittleEndian, buffer[0:20])
	header = append(header, buf.Bytes()...)
	buf.Reset()

	header = append(header, []byte("data")...)
	binary.Write(buf, binary.LittleEndian, length)
	header = append(header, buf.Bytes()...)
	return header
}
