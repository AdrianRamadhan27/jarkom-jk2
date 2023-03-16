package main

import (
	"net"
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

func main() {
	//The Program logic should go here.

}

func Fetch(req HttpRequest, connection net.Conn) (HttpResponse, []Student, HttpRequest) {
	//This program handles the request-making to the server
	var res HttpResponse
	var Student []Student

	return res, Student, req

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
		} else {
			str = str + string(stringByte)
		}
		res.Data = str

	}
	return res

}

func RequestEncoder(req HttpRequest) []byte {
	var result string
	result = req.Method + " " + req.Uri + " " + req.Version + "\r\n" + req.Host + "\r\n" + req.Accept + "\r\n" + req.AcceptLanguange
	return []byte(result)

}
