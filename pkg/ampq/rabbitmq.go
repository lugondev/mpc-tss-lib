package rabbitmq

import (
	"context"
	"github.com/lugondev/mpc-tss-lib/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"time"
)

type BridgeMQ struct {
	config       *config.RabbitConfig
	exchangeName string
	topic        string

	conn     *amqp.Connection
	ch       *amqp.Channel
	q        *amqp.Queue
	ctx      context.Context
	cancel   context.CancelFunc
	MsgCh    <-chan amqp.Delivery
	confirms chan amqp.Confirmation
	errCh    chan error
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal().Err(err).Msg(msg)
	}
}

func InitMqClient(rabbitConfig *config.RabbitConfig, exchange, topic string) *BridgeMQ {
	//confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(rabbitConfig.Timeout)*time.Second)

	mq := &BridgeMQ{
		config: rabbitConfig,
		q:      nil,
		ctx:    ctx,
		cancel: cancel,
		//confirms: confirms,
		exchangeName: exchange,
		topic:        topic,
	}
	mq.Connect()
	return mq
}

func (mq *BridgeMQ) ConsumerWarrantConnection() {
	for {
		select {
		case err := <-mq.errCh:
			if err != nil {
				mq.Connect()
				mq.Subscribe(false)
			}
		}
	}
}

func (mq *BridgeMQ) Connect() {
	var err error
	log.Info().Msgf("Connecting to mq: %s", mq.config.Url)
	mq.conn, err = amqp.Dial(mq.config.Url)
	failOnError(err, "Failed to open connection to RabbitMQ")

	mq.ch, err = mq.conn.Channel()
	failOnError(err, "Failed to open channel")
	go func() {
		<-mq.conn.NotifyClose(make(chan *amqp.Error)) //Listen to NotifyClose
		log.Error().Msg("connection closed")
	}()
	err = mq.ch.ExchangeDeclare(
		mq.exchangeName, // name
		"fanout",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to init Exchange")
	q, err := mq.ch.QueueDeclare(
		mq.topic, // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = mq.ch.QueueBind(
		q.Name,          // queue name
		"",              // routing key
		mq.exchangeName, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")
	//confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))
}

func (mq *BridgeMQ) Subscribe(autoAck bool) {
	var err error
	mq.MsgCh, err = mq.ch.Consume(
		mq.topic, // queue
		"",       // consumer
		autoAck,  // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	failOnError(err, "Failed to register a consumer")
	log.Info().Msg("Subscribe topic success")
	go mq.ConsumerWarrantConnection()
}

func (mq *BridgeMQ) Publish(msg []byte, id string) {
	select {
	case err := <-mq.errCh:
		if err != nil {
			mq.Connect()
		}
	default:

	}
	log.Info().Bytes(id, msg).Msg("Publish message")
	defer mq.cancel()
	err := mq.ch.PublishWithContext(
		mq.ctx,
		mq.exchangeName,
		"", false, false,
		amqp.Publishing{
			DeliveryMode: 2,
			ContentType:  "text/plain",
			Body:         msg,
			MessageId:    id,
		},
	)
	failOnError(err, "Failed to publish a message")
	log.Info().Msgf("Published to mq: id: %s msg: %s", id, msg)
}
