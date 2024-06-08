package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"

	"github.com/gofiber/websocket/v2"
)

type WebsocketServer struct {
	clients   map[*websocket.Conn]bool
	broadcast chan *Message
}

func NewWebSocket() *WebsocketServer {
	return &WebsocketServer{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan *Message),
	}
}

func (s *WebsocketServer) HandleWebSocket(ctx *websocket.Conn) {
	s.clients[ctx] = true
	defer func() {
		delete(s.clients, ctx)
		ctx.Close()
	}()

	for {
		_, msg, err := ctx.ReadMessage()
		if err != nil {
			log.Println("Read Error: ", err)
			break
		}

		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Fatalf("Error Unmarshlling")
		}
		s.broadcast <- &message
	}
}

func (s *WebsocketServer) HandleMessages() {
	for {
		msg := <-s.broadcast

		for client := range s.clients {
			err := client.WriteMessage(websocket.TextMessage, getMessageTemplate(msg))
			if err != nil {
				log.Println("Write Error: %v", err)
				client.Close()
				delete(s.clients, client)
			}
		}
	}
}

func getMessageTemplate(msg *Message) []byte {
	tmpl, err := template.ParseFiles("views/message.html")
	if err != nil {
		log.Fatalf("template parsing failed: %s", err)
	}
	var renderedMsg bytes.Buffer
	err = tmpl.Execute(&renderedMsg, msg)
	if err != nil {
		log.Fatalf("template exec failed: %s", err)
	}
	return renderedMsg.Bytes()
}
