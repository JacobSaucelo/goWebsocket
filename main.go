package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.Handle("/order-book", websocket.Handler(server.handleWsOrderBook))
	fmt.Println("Running on localhost:3000")
	http.ListenAndServe(":3000", nil)
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWsOrderBook(ws *websocket.Conn) {
	fmt.Println("incomming connection from client to OrderBook feed: ", ws.RemoteAddr())

	for {
		payload := fmt.Sprintf("orderbook data -> %d\n", time.Now().UnixNano())
		ws.Write([]byte(payload))
		time.Sleep(time.Second * 2)
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("incomming connection from client:", ws.RemoteAddr())

	s.conns[ws] = true
	// mutex
	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				// end of file
				break
			}
			fmt.Println("Read error: ", err)
			//! DO NOT BREAK, CONNECTION WILL GET CUT
			continue
		}

		msg := buf[:n]
		s.broadcast(msg)
		// fmt.Println(string(msg))
		// ws.Write([]byte("CHIPI CHIPI CHAPA CHAPA DUBI DUBI!!"))
	}
}

func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error: ", err)
			}
		}(ws)
	}
}

/*
	for future projects js

	let socket = new WebSocket("ws://localhost:3000/ws")
	socket.onmessage = (e) => { console.log("received message: ", e.data)}
	socket.send("21 can you do sum for me")
*/
