package websocket

import (
	"api2/src/models/services"
	"api2/utils"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]string)
	broadcast = make(chan []byte)
	mutex     = &sync.Mutex{}
)

const (
	pingInterval     = 30 * time.Second
	pongWait         = 60 * time.Second
	writeWait        = 10 * time.Second
)

func HandleConnections(c *gin.Context) {
	claims, err := utils.ValidateTokenFromQuery(c, "admin", "user")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
		return
	}

	log.Printf("✅ Usuario autenticado: ID=%d, Rol=%s, Zona=%s", claims.UserID, claims.Role, claims.Zona)

	upgrader := websocket.Upgrader{
		CheckOrigin:      func(r *http.Request) bool { return true },
		HandshakeTimeout: 30 * time.Second,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("❌ WebSocket upgrade error: %v", err)
		return
	}

	defer func() {
		utils.RemoveClient(ws)
		ws.Close()
		log.Println("🔴 Conexión WebSocket cerrada")
	}()

	utils.RegisterClient(ws, claims.Zona)
	go services.StartDynamicConsumerByZona(claims.Zona)

	// Set timeout para espera de pong
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// 🟡 Rutina para enviar pings cada cierto tiempo
	go func(conn *websocket.Conn) {
		ticker := time.NewTicker(pingInterval)
		defer ticker.Stop()

		for range ticker.C {
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("⚠️ Error al enviar ping: %v", err)
				conn.Close()
				return
			}
		}
	}(ws)

	log.Println("🟢 WebSocket activo.")

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("🔴 Conexión cerrada por error o timeout: %v", err)
			break
		}
		log.Printf("📩 Mensaje recibido: %s", msg)
	}
}

func NotifyClients(data any) {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println("Error al serializar data para broadcast:", err)
		return
	}
	log.Printf("📤 Enviando mensaje al canal broadcast: %s", string(bytes))
	broadcast <- bytes
}

func StartBroadcaster() {
	for {
		msg := <-broadcast
		log.Printf("📡 Broadcast: enviando mensaje a %d clientes", len(clients))
		mutex.Lock()
		for client := range clients {
			client.SetWriteDeadline(time.Now().Add(writeWait))
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("Error al enviar mensaje: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}


