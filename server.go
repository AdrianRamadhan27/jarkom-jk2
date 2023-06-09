package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net"
	"os"
)

const (
	SERVER_HOST = ""
	SERVER_PORT = "2376"
	SERVER_TYPE = "tcp"
	BUFFER_SIZE = 1024
	GROUP_NAME  = "JK2"
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

type Student struct {
	Nama string
	Npm  string
}

func main() {
	//The Program logic should go here.
	tcpAddr, err := net.ResolveTCPAddr(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error message:", err.Error())
	}

	server, err := net.ListenTCP(SERVER_TYPE, tcpAddr)
	if err != nil {
		fmt.Println("Error message:", err.Error())
		os.Exit(1)
	}

	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error message: ", err.Error())
		}

		fmt.Println("Accept connection from: ", conn.RemoteAddr().String())
		go HandleConnection(conn)
	}
}

func HandleConnection(connection net.Conn) {
	//This progrom handles the incoming request from client
	buffer := make([]byte, BUFFER_SIZE)
	bufLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error message:", err.Error())
		return
	}

	reqByte := buffer[:bufLen]

	req := RequestDecoder(reqByte)

	fmt.Printf("Received %s request from %s: %s\n", req.Method, connection.RemoteAddr().String(), req.Uri)

	resp := HandleRequest(req)

	respByte := ResponseEncoder(resp)

	_, err = connection.Write(respByte)
	if err != nil {
		fmt.Println("error message:", err.Error())
		return
	}

	defer connection.Close()
}

func HandleRequest(req HttpRequest) HttpResponse {
	//This program handles the routing to each view handler.
	var status string
	var data string
	var contentType string
	var contentLanguage string
	if req.Uri == "/" || req.Uri == "/?name="+GROUP_NAME {
		status = "200"
		data = "<html><body><h1>Halo, kami dari Klepon</h1></body></html>"
		contentType = "text/html"
		contentLanguage = req.AcceptLanguange
	} else if req.Uri == "/data" {
		status = "200"
		var students [3]Student
		students[0] = Student{Nama: "Raden Mohamad Adrian Ramadhan Hendar Wibawa", Npm: "2106750540"}
		students[1] = Student{Nama: "Hizkia Sebastian Ginting", Npm: "2106750881"}
		students[2] = Student{Nama: "Kade Satrya Noto Sadharma", Npm: "2106752376"}
		switch req.Accept {
		case "application/json":
			jsonData, err := json.Marshal(students)
			if err != nil {
				data = "Error message: " + err.Error()
			}
			data = string(jsonData)
			contentType = "application/json"
		case "application/xml":
			xmlData, err := xml.Marshal(students)
			if err != nil {
				data = "Error message: " + err.Error()
			}
			data = string(xmlData)
			contentType = "application/xml"
		default:
			jsonData, err := json.Marshal(students)
			if err != nil {
				data = "Error message: " + err.Error()
			}
			data = string(jsonData)
			contentType = "application/xml"
		}
		contentLanguage = req.AcceptLanguange
	} else if req.Uri == "/greeting" {
		status = "200"
		switch req.AcceptLanguange {
		case "id-ID":
			data = "<html><body><h1>Halo, kami dari Klepon</h1></body></html>"
			contentLanguage = "id-ID"
		case "en-US":
			data = "<html><body><h1>Hello, we are from Klepon</h1></body></html>"
			contentLanguage = "en-US"
		default:
			data = "<html><body><h1>Hello, we are from Klepon</h1></body></html>"
			contentLanguage = "en-US"
		}
		contentType = "text/html"
	} else {
		status = "404"
		data = ""
		contentType = "text/html"
		contentLanguage = req.AcceptLanguange
	}

	return HttpResponse{
		Version:         req.Version,
		StatusCode:      status,
		ContentType:     contentType,
		ContentLanguage: contentLanguage,
		Data:            data,
	}

}

func RequestDecoder(bytestream []byte) HttpRequest {
	//Put the decoding program for HTTP Request Packet here
	var req HttpRequest
	var stringByte byte
	var loopControl int = 0
	var skippedIndex int = -1
	var str string = ""

	for i := 0; i < len(bytestream); i++ {
		stringByte = bytestream[i]
		if i == skippedIndex {
			continue
		}

		if stringByte == 32 && loopControl <= 1 {
			switch loopControl {
			case 0:
				req.Method = str
			case 1:
				req.Uri = str
			}
			str = ""
			loopControl++
		} else if stringByte == 13 && bytestream[i+1] == 10 && loopControl > 1 {
			switch loopControl {
			case 2:
				req.Version = str
			case 3:
				req.Host = str
			case 4:
				req.Accept = str
			case 5:
				req.AcceptLanguange = str
			}
			str = ""
			skippedIndex = i + 1
			loopControl++
		} else {
			str = str + string(stringByte)
		}

	}
	return req
}

func ResponseEncoder(res HttpResponse) []byte {
	//Put the encoding program for HTTP Response Struct here
	var result string
	result = res.Version + " " + res.StatusCode + "\r\n" + res.ContentType + "\r\n" + res.ContentLanguage + "\r\n\r\n" + res.Data
	return []byte(result)
}
