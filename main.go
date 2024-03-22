package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type request struct {
	Cmd         string `json:"cmd"`
	ChannelName string `json:"chName"`
	Msg         string `json:"msg"`
}

var mutex sync.Mutex

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Update with your Redis server address
		Password: "",               // No password by default
		DB:       0,                // Default DB
	})
}

func handleWebSocket(c *gin.Context) {
	userId := c.Query("userId")
	if userId == "" {
		log.Println("userId is empty")
		c.JSON(http.StatusBadRequest, gin.H{"err": "userId not provided"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// maintaining list of all subscribed pubsub objects, so once connection is disconnected, we can close all pubsub object
	var subscribedChannels = make(map[string]*redis.PubSub)

	defer func() {
		log.Printf("closing total pubsubs: %d\n", len(subscribedChannels))
		for _, ps := range subscribedChannels {
			ps.Close()
		}
		log.Println("closing connection")
		conn.Close()
	}()

	for {
		// Read message from the client
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var data request
		err = json.Unmarshal(msg, &data)
		if err != nil {
			continue
		}
		switch data.Cmd {
		case "I-JC": // Join Channel
			// if not already subscribed, then only subscribe
			if _, ok := subscribedChannels[data.ChannelName]; !ok {
				pubsub := redisClient.Subscribe(c, data.ChannelName)
				go listenToChannel(conn, pubsub)
				subscribedChannels[data.ChannelName] = pubsub
			}
		case "I-SM": // Send Message
			err = redisClient.Publish(c, data.ChannelName, userId+": "+data.Msg).Err()
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// listen for redis messages in a goroutine
func listenToChannel(conn *websocket.Conn, ps *redis.PubSub) {
	for {
		msg, ok := <-ps.Channel()
		if !ok {
			return
		}
		mutex.Lock()
		conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
		mutex.Unlock()
	}
}

func main() {
	r := gin.Default()

	r.GET("/chat", handleWebSocket)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
