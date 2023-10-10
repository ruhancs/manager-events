package main

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	"github.com/ruhancs/manager-events/pkg/rabbitmq"
)

func main() {
	ch,err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	
	msgsOut := make(chan amqp091.Delivery)

	go rabbitmq.Consumer(ch,msgsOut,"test_queue")
	
	for msg := range msgsOut {
		fmt.Println(string(msg.Body))
		//apagar msg da fila
		msg.Ack(false)
	}
}