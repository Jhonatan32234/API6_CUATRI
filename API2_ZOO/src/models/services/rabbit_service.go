package services

import (
	"api2/src/entities"
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func StartRabbitConsumers() {
	url := os.Getenv("RABBITCONN")
	conn, err := amqp.Dial(url)

	if err != nil {
		log.Fatal("RabbitMQ connection failed:", err)
	}
	ch, _ := conn.Channel()

	go consumeTopic(ch, "visitas_topic", "visita.data", handleVisitas)
	go consumeTopic(ch, "atracciones_topic", "atraccion.data", handleAtracciones)
}

func consumeTopic(ch *amqp.Channel, exchange, routingKey string, handler func([]byte)) {
	err := ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("❌ Error al declarar el exchange:", err)
		return
	}

	q, err := ch.QueueDeclare(
		routingKey+"_queue", 
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("❌ Error al declarar la cola:", err)
		return
	}

	err = ch.QueueBind(
		q.Name,
		routingKey,
		exchange,
		false,
		nil,
	)
	if err != nil {
		log.Println("❌ Error al enlazar la cola al tópico:", err)
		return
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Println("❌ Error al iniciar el consumo:", err)
		return
	}

	go func() {
		for d := range msgs {
			log.Println("📨 Mensaje recibido:", string(d.Body))
			handler(d.Body)
			ch.Ack(d.DeliveryTag, false)
		}
	}()
}

// Procesador de visitas
func handleVisitas(data []byte) {
	var visitas []entities.Visitas
	if err := json.Unmarshal(data, &visitas); err != nil {
		log.Println("❌ Error al deserializar visitas:", err)
		return
	}

	guardadas, err := SaveVisitas(visitas)
	if err != nil {
		log.Println("❌ Error al guardar visitas:", err)
		return
	}

	log.Printf("✅ Visitas guardadas: %d", len(guardadas))
}

// Procesador de atracciones
func handleAtracciones(data []byte) {
	var atracciones []entities.Atraccion
	if err := json.Unmarshal(data, &atracciones); err != nil {
		log.Println("❌ Error al deserializar atracciones:", err)
		return
	}

	guardadas, err := SaveAtracciones(atracciones)
	if err != nil {
		log.Println("❌ Error al guardar atracciones:", err)
		return
	}
	log.Printf("✅ Atracciones guardadas: %d", len(guardadas))
}
