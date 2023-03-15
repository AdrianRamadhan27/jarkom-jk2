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

	return res

}

func RequestEncoder(req HttpRequest) []byte {
	var result string

	return []byte(result)

}
