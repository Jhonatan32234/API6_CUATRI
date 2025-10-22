package utils

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]string)
	Broadcast = make(chan []byte)
	Mutex     = &sync.Mutex{}
)

func RegisterClient(ws *websocket.Conn, zona string) {
	Mutex.Lock()
	defer Mutex.Unlock()
	clients[ws] = zona
	log.Printf("👤 Cliente registrado para zona: %s", zona)
}

func RemoveClient(ws *websocket.Conn) {
	Mutex.Lock()
	defer Mutex.Unlock()
	delete(clients, ws)
}

func NotifyClients(data map[string]interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println("❌ Error al serializar datos para WebSocket:", err)
		return
	}

	// Leer zona directamente del mapa principal
	zona, ok := data["zona"].(string)
	if !ok || zona == "" {
		log.Println("⚠️ Zona no encontrada en los datos. Broadcast cancelado.")
		return
	}

	log.Printf("📡 Broadcast activado para zona: %s", zona)

	Mutex.Lock()
	defer Mutex.Unlock()

	for client, clientZona := range clients {
		if clientZona == zona {
			err := client.WriteMessage(websocket.TextMessage, bytes)
			if err != nil {
				log.Printf("❌ Error al enviar a cliente zona '%s': %v", clientZona, err)
				client.Close()
				delete(clients, client)
			} else {
				log.Println("✅ Mensaje enviado a cliente de zona:", clientZona)
			}
		}
	}
}



func StartBroadcaster() {
	for {
		msg := <-Broadcast
		log.Printf("📡 Broadcast activado: %s\n", string(msg))

		Mutex.Lock()
		if len(clients) == 0 {
			log.Println("⚠️ No hay clientes WebSocket conectados para recibir el mensaje.")
		}

		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("❌ Error al enviar mensaje a un cliente: %v", err)
				client.Close()
				delete(clients, client)
			} else {
				log.Println("✅ Mensaje enviado correctamente a un cliente.")
			}
		}
		Mutex.Unlock()
	}
}

