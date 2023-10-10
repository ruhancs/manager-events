package main

import "github.com/ruhancs/manager-events/pkg/rabbitmq"

func main() {
	ch,err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	rabbitmq.Publish(ch, "Testando...", "amq.direct")
	rabbitmq.Publish(ch, "Testando2...", "amq.direct")
	rabbitmq.Publish(ch, "Testando3...", "amq.direct")
	rabbitmq.Publish(ch, "Testando4...", "amq.direct")
}