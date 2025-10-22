package services

import (
	"api2/utils"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/streadway/amqp"
)


var (
	zonaConsumers     = make(map[string]bool)
	zonaConsumersLock = &sync.Mutex{}
)

func StartDynamicConsumerByZona(zona string) {
	zonaConsumersLock.Lock()
	defer zonaConsumersLock.Unlock()

	if zonaConsumers[zona] {
		log.Printf("⚠️ Consumidor ya iniciado para la zona: %s", zona)
		return
	}

	zonaConsumers[zona] = true
	log.Printf("🚀 Iniciando consumidores para la zona: %s", zona)

	go consumeZonaTopic("visitas_topic", fmt.Sprintf("visitas.%s", zona),
	func(id uint) {
		handleZonaVisita(id, zona)
	})

go consumeZonaTopic("atracciones_topic", fmt.Sprintf("atracciones.%s", zona),
	func(id uint) {
		handleZonaAtraccion(id, zona)
	})

}


func consumeZonaTopic(exchange, routingKey string, handler func(uint)) {
	url := os.Getenv("RABBITCONN")
	log.Printf("📡 Iniciando consumidor para zona: exchange='%s', routingKey='%s'\n", exchange, routingKey)
	conn, err := amqp.Dial(url)

	if err != nil {
		log.Println("❌ RabbitMQ conexión fallida:", err)
		return
	}
	log.Println("✅ Conectado a RabbitMQ para zona.")
	ch, _ := conn.Channel()

	err = ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	if err != nil {
		log.Println("❌ Error declarando exchange:", err)
		return
	}

	q, err := ch.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		log.Println("❌ Error declarando cola:", err)
		return
	}

	err = ch.QueueBind(q.Name, routingKey, exchange, false, nil)
	if err != nil {
		log.Println("❌ Error enlazando cola a tópico:", err)
		return
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Println("❌ Error al consumir:", err)
		return
	}

	go func() {
		log.Println("🟢 Esperando mensajes de zona en:", routingKey)
		for d := range msgs {
			log.Printf("📥 Mensaje recibido de zona (%s): %s\n", routingKey, string(d.Body))

			var payload struct {
				Id uint `json:"id"`
			}
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				log.Println("❌ Error parsing ID payload:", err)
				continue
			}

			log.Printf("🔍 ID extraído del mensaje: %d\n", payload.Id)

			go func(id uint) {
				time.Sleep(1 * time.Second)
				handler(id)
			}(payload.Id)
		}
	}()
}

func handleZonaVisita(id uint, zona string) {
	log.Printf("📌 Procesando visita por zona con ID: %d (zona: %s)\n", id, zona)

	utils.NotifyClients(map[string]interface{}{
		"type": "visita",
		"zona": zona,
	})
	log.Println("✅ Enviado al WebSocket.")
}

func handleZonaAtraccion(id uint, zona string) {
	log.Printf("📌 Procesando atracción por zona con ID: %d (zona: %s)\n", id, zona)

	utils.NotifyClients(map[string]interface{}{
		"type": "atraccion",
		"zona": zona,
	})
	log.Println("✅ Enviado al WebSocket.")
}


