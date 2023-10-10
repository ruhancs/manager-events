package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

// criar conexao e canal no rabbitmq
func OpenChannel() (*amqp.Channel, error) {
	//chamada para criar a conexao, namespace é /
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	//criar channel do rabbitmq
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return ch, nil
}

// consumidor de msg do rabbitmq, msgsOut chan amqp.Delivery tem todas as informacoes das menssagens recebidas
func Consumer(ch *amqp.Channel, msgOut chan amqp.Delivery, queue string) error {
	//consumir as informacoes da fila do rabbitmq
	msgs, err := ch.Consume(
		queue,  //nome da fila que sera consumida
		"go-consumer", //nome da aplicacao que ira consumir
		false,         //autoAck apagar automaticamente as messages
		false,         //fila exclusiva
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	//ler as menssagens recebidas na fila
	for msg := range msgs {
		//joga as mensagens recebidas do rabbitmq no canal msgOut
		msgOut <- msg
	}

	return nil
}

// publicar menssagens no rabbitmq
func Publish(ch *amqp.Channel, body string, exchange string) error {
	err := ch.PublishWithContext(
		context.Background(),
		//exchange que sera enviada a msg,exchange padrao do rabbitmq amq.direct
		exchange,
		//key, para a exchange mandar as msgs para fila que possui a key
		"",
		false,
		false,
		//como sera a msg enviada, ex text/plain, json
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body), // convert string para []byte
		},
	)
	if err != nil {
		return err
	}

	return nil
}
