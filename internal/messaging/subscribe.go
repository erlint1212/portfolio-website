package messaging

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/erlint1212/portfolio/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func WriteLog(gl routing.GameLog) error {
	log.Printf("[GameLog] %s: %s", gl.CurrentTime.Format("15:04:05"), gl.Message)
	return nil
}

func UnmarshalGob[T any](data []byte) (T, error) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	var target T
	err := decoder.Decode(&target)
	return target, err
}

func HandlerWriteLog() func(gl routing.GameLog) routing.AckType {
	return func(gl routing.GameLog) routing.AckType {
		err := WriteLog(gl)
		if err != nil {
			log.Printf("[ERROR] Couldn't write to log: %v\n", err)
			return routing.NackDiscard
		}
		return routing.Ack
	}
}

func Subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType routing.SimpleQueueType,
	handler func(T) routing.AckType,
	unmarshaller func([]byte) (T, error),
) error {
	ch, queue, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return fmt.Errorf("could not declare and bind queue: %v", err)
	}

	err = ch.Qos(10, 0, false)
	if err != nil {
		return fmt.Errorf("could not set qos: %v", err)
	}

	msgs, err := ch.Consume(
		queue.Name, "", false, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("could not consume messages: %v", err)
	}

	go func() {
		defer ch.Close()
		for msg := range msgs {
			target, err := unmarshaller(msg.Body)
			if err != nil {
				fmt.Printf("could not unmarshal message: %v\n", err)
				msg.Ack(false) 
				continue
			}
			acktype := handler(target)
			switch acktype {
			case routing.Ack:
				msg.Ack(false)
			case routing.NackRequeue:
				msg.Nack(false, true)
			case routing.NackDiscard:
				msg.Nack(false, false)
			}
		}
	}()
	return nil
}
