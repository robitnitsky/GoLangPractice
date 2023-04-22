package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var conns = map[string]*websocket.Conn{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	e := gin.Default()
	e.GET("/api/v1/message", SendMessage)
	panic(http.ListenAndServe(":8080", e))
}

func SendMessage(c *gin.Context) {
	var reqFrom string
	h := http.Header{}

	for _, sub := range websocket.Subprotocols(c.Request) {
		h.Set("Sec-Websocket-Protocol", sub)
		reqFrom = sub
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, h)
	if err != nil {
		log.Println(err)
	}

	conns[reqFrom] = ws

	for {
		var req Request

		err = ws.ReadJSON(&req)
		if err != nil {
			log.Println(err)
		}

		if con, ok := conns[req.To]; ok {
			err = con.WriteJSON(&req)
			if err != nil {
				log.Println(err)
			}
		}
	}

}

type Request struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
}
