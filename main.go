package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var redisClient *redis.Client

const ChatChannel = "chat"

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Update with your Redis server address
		Password: "",               // No password by default
		DB:       0,                // Default DB
	})
}

func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	pubsub := redisClient.Subscribe(c, ChatChannel)
	defer pubsub.Close()

	//listen for redis messages in a goroutine
	msgChannel := pubsub.Channel()
	go func() {
		for {
			msg, ok := <-msgChannel
			if !ok {
				return
			}
			conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
		}
	}()

	for {
		// Read message from the client
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if msgType == websocket.TextMessage {
			err := redisClient.Publish(c, ChatChannel, msg).Err()
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func main() {
	r := gin.Default()

	r.GET("/ws", handleWebSocket)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
