package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type HttpRequest struct {
	Method          string
	Uri             string
	Version         string
	Host            string
	Accept          string
	AcceptLanguange string
}

type HttpResponse struct {
	Version         string
	StatusCode      string
	ContentType     string
	ContentLanguage string
	Data            string
}

const (
	SERVER_TYPE = "tcp"
	BUFFER_SIZE = 1024
)

type Student struct {
	Nama string
	Npm  string
}

func Split(r rune) bool {
	return r == ':' || r == '/'
}

func main() {
	//The Program logic should go here.
	var url string
	var data string
	var lang string
	fmt.Print("input the url: ")
	fmt.Scan(&url)
	a := strings.FieldsFunc(url, Split)
	fmt.Print("input the data type: ")
	fmt.Scan(&data)
	fmt.Print("input the language: ")
	fmt.Scan(&lang)
	var uri string
	if len(a) == 3 {
		uri = "/"
	} else {
		uri = "/" + a[3]
	}
	var req HttpRequest
	req = HttpRequest{
		Method:          "GET",
		Uri:             uri,
		Version:         "HTTP/1.1",
		Host:            a[1],
		Accept:          data,
		AcceptLanguange: lang,
	}
	tcpAddr, err := net.ResolveTCPAddr(SERVER_TYPE, req.Host+":"+a[2])
	if err != nil {
		fmt.Println("Error message:", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP(SERVER_TYPE, nil, tcpAddr)
	if err != nil {
		fmt.Println("Error message:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	var res HttpResponse
	var student []Student
	res, student, req = Fetch(req, conn)
	_ = res
	_ = student
	fmt.Printf("Status code: %s\n", res.StatusCode)
	fmt.Printf("Body: %s", res.Data)
}

func Fetch(req HttpRequest, connection net.Conn) (HttpResponse, []Student, HttpRequest) {
	//This program handles the request-making to the server
	var res HttpResponse
	var student []Student
	var reqData []byte

	reqData = RequestEncoder(req)

	_, err := connection.Write(reqData)
	if err != nil {
		fmt.Println("Error message:", err.Error())
		os.Exit(1)
	}

	buffer := make([]byte, BUFFER_SIZE)
	bufLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error message:", err.Error())
		os.Exit(1)
	}

	res = ResponseDecoder(buffer[:bufLen])
	return res, student, req
}

func ResponseDecoder(bytestream []byte) HttpResponse {
	var res HttpResponse
	var stringByte byte
	var loopControl int = 0
	var skippedIndex int = -1
	var str string = ""

	for i := 0; i < len(bytestream); i++ {
		stringByte = bytestream[i]
		if i == skippedIndex {
			continue
		}

		if stringByte == 32 && loopControl == 0 {
			res.Version = str
			str = ""
			loopControl++
		} else if stringByte == 13 && bytestream[i+1] == 10 && loopControl > 0 {
			switch loopControl {
			case 1:
				res.StatusCode = str
			case 2:
				res.ContentType = str
			case 3:
				res.ContentLanguage = str
			}
			str = ""
			skippedIndex = i + 1
			loopControl++
		} else {
			str = str + string(stringByte)
		}
		res.Data = str

	}
	return res

}

func RequestEncoder(req HttpRequest) []byte {
	var result string
	result = req.Method + " " + req.Uri + " " + req.Version + "\r\n" + req.Host + "\r\n" + req.Accept + "\r\n" + req.AcceptLanguange + "\r\n\r\n"
	return []byte(result)

}
