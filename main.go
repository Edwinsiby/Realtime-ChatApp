package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	connectedPublicClients = make(map[string]*websocket.Conn)
	userName               string
	err                    error
)

func main() {
	router := gin.New()
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/static", "./static")
	router.GET("/", Index)
	router.POST("user-login", Login)
	router.GET("ws/chat", Chat)
	router.Run()
}

func Index(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func Login(c *gin.Context) {
	userName = c.Request.FormValue("userName")
	if _, ok := connectedPublicClients[userName]; ok {
		c.HTML(http.StatusBadRequest, "index.html", "user already exsist")
	}
	c.Set("userName", userName)
	c.HTML(200, "chat.html", userName)
}

func Chat(c *gin.Context) {
	userName = c.GetString("userName")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("socket connection err :", err)
		return
	}
	defer func() {
		delete(connectedPublicClients, userName)
		conn.Close()
	}()
	if _, ok := connectedPublicClients[userName]; !ok {
		connectedMessage := []byte("Connected to the server")
		err = conn.WriteMessage(websocket.TextMessage, connectedMessage)
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
		connectedPublicClients[userName] = conn
	}
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}
		fmt.Println(p)
		for clientName, clientConn := range connectedPublicClients {
			if clientName != userName {
				err = clientConn.WriteMessage(messageType, p)
				if err != nil {
					log.Println("Error forwarding message:", err)
					return
				}
			}
		}
	}
}
