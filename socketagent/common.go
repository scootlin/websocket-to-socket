package socketagent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

var funcRoute = map[uint16]func(Packet){
	// UDP
	4354: func4354,
	4609: func4609,
	4864: func4864,
	5377: func5377,
	9730: func9730,
	9731: func9731,
	// TCP
	5121: func5121,
	5123: func5121,
}

func doFunc(datatype uint16, pack Packet) {
	if fn, ok := funcRoute[datatype]; ok {
		fn(pack)
	}
}

func setCalltype(calltype string) uint8 {
	n, _ := strconv.Atoi(calltype)
	return uint8(n)
}

func setCallee(callee string) uint32 {
	n, _ := strconv.Atoi(callee)
	return uint32(n)
}

func setCaller(caller string) uint32 {
	n, _ := strconv.Atoi(caller)
	return uint32(n)
}

func hexdec(datatype uint16) uint16 {
	s := strconv.FormatUint(uint64(datatype), 10)
	s = "0x" + s
	n, _ := strconv.ParseUint(s, 0, 32)
	return uint16(n)
}

func createHeader(opcode byte, calltype byte, callee uint32, caller uint32, datatype uint16) *Header {
	header := &Header{}
	header.Version = 1
	header.Opcode = opcode
	header.Calltype = calltype
	header.Callee = make([]byte, 16)
	header.Caller = make([]byte, 16)
	header.Session = 0
	header.Datatype = datatype
	header.Protocol = 8
	header.Length = 0
	binary.LittleEndian.PutUint32(header.Callee, callee)
	binary.LittleEndian.PutUint32(header.Caller, caller)
	return header
}

func encodeToBytes(buf *bytes.Buffer, data interface{}) {
	var p reflect.Value
	switch reflect.ValueOf(data).Kind() {
	case reflect.Ptr:
		p = reflect.ValueOf(data).Elem()
	default:
		p = reflect.ValueOf(data)
	}

	if p.Kind() == reflect.Struct {
		for i := 0; i < p.NumField(); i++ {
			v := p.Field(i)
			if fmt.Sprint(v) != "<nil>" {
				if v.Kind() == reflect.Ptr && v.Type().Elem().Kind() == reflect.Struct {
					encodeToBytes(buf, v.Interface())
				} else if v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.Struct {
					for j := 0; j < v.Len(); j++ {
						encodeToBytes(buf, v.Index(j).Interface())
					}
				} else {
					binary.Write(buf, binary.LittleEndian, v.Interface())
				}
			}
		}
	} else if p.Kind() == reflect.Slice {
		for i := 0; i < p.Len(); i++ {
			encodeToBytes(buf, p.Index(i).Interface())
		}
	}
}

func unpackSignedInteger32(data []byte) int32 {
	var p int32
	buf := bytes.NewReader(data)
	binary.Read(buf, binary.LittleEndian, &p)
	return p
}

func unpackSignedInteger16(data []byte) int16 {
	var p int16
	buf := bytes.NewReader(data)
	binary.Read(buf, binary.LittleEndian, &p)
	return p
}

func CheckError(err error) {
	if err != nil {
		log.Println(err)
		return
		// os.Exit(0)
		// panic(fmt.Sprintln(err))
	}
}
