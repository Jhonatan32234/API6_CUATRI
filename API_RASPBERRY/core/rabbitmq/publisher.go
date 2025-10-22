package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

func PublishToTopic[T any](data []T, exchangeName string, routingKey string) bool {
	log.Println("🚀 Iniciando PublishToTopic...")

	log.Println("🔌 Intentando conectar a RabbitMQ...")
	conn, err := amqp.Dial("amqp://admin:password@localhost:5672/")
	if err != nil {
		log.Println("❌ Falló la conexión a RabbitMQ:", err)
		return false
	}
	log.Println("✅ Conexión establecida.")
	defer func() {
		log.Println("🔌 Cerrando conexión a RabbitMQ...")
		conn.Close()
	}()

	log.Println("📡 Creando canal...")
	ch, err := conn.Channel()
	if err != nil {
		log.Println("❌ Falló al abrir un canal:", err)
		return false
	}
	log.Println("✅ Canal creado.")
	defer func() {
		log.Println("📡 Cerrando canal...")
		ch.Close()
	}()

	log.Printf("📦 Declarando exchange '%s' tipo 'topic'...\n", exchangeName)
	err = ch.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("❌ Falló al declarar el exchange:", err)
		return false
	}
	log.Println("✅ Exchange declarado.")

	log.Println("🧬 Serializando datos a JSON...")
	body, err := json.Marshal(data)
	if err != nil {
		log.Println("❌ Error al serializar los datos:", err)
		return false
	}
	log.Printf("📨 Publicando mensaje en el exchange '%s' con routingKey '%s'...\n", exchangeName, routingKey)
	err = ch.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Println("❌ Falló al publicar:", err)
		return false
	}
	log.Println("✅ Publicado correctamente en el topic:", routingKey)
	return true
}

func PublishIDToZoneTopic(exchangeName string, zona string, id int, tipo string) bool {
	log.Println("🚀 Iniciando PublishIDToZoneTopic...")

	log.Println("🔌 Intentando conectar a RabbitMQ...")
	conn, err := amqp.Dial("amqp://admin:password@localhost:5672/")
	if err != nil {
		log.Println("❌ Falló la conexión a RabbitMQ:", err)
		return false
	}
	log.Println("✅ Conexión establecida.")
	defer func() {
		log.Println("🔌 Cerrando conexión a RabbitMQ...")
		conn.Close()
	}()

	log.Println("📡 Creando canal...")
	ch, err := conn.Channel()
	if err != nil {
		log.Println("❌ Falló al abrir un canal:", err)
		return false
	}
	log.Println("✅ Canal creado.")
	defer func() {
		log.Println("📡 Cerrando canal...")
		ch.Close()
	}()

	log.Printf("📦 Declarando exchange '%s' tipo 'topic'...\n", exchangeName)
	err = ch.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("❌ Falló al declarar el exchange:", err)
		return false
	}
	log.Println("✅ Exchange declarado.")

	routingKey := tipo + "." + zona
	log.Printf("🔧 Generando routingKey: %s\n", routingKey)

	log.Printf("🧬 Serializando ID (%d) a JSON...\n", id)
	
	body, err := json.Marshal(map[string]int{"id": id})
	if err != nil {
		log.Println("❌ Error al serializar el ID:", err)
		return false
	}
	log.Printf("📨 Payload JSON: %s\n", string(body))

	log.Printf("📨 Publicando mensaje en el exchange '%s' con routingKey '%s'...\n", exchangeName, routingKey)
	err = ch.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Println("❌ Falló al publicar en zona:", err)
		return false
	}
	log.Println("✅ Publicado en zona:", routingKey)
	return true
}
