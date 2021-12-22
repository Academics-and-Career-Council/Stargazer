package Services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Academics-and-Career-Council/Stargazer.git/internal/api"
	"github.com/Academics-and-Career-Council/Stargazer.git/internal/database"
	"github.com/Academics-and-Career-Council/Stargazer.git/internal/models"
	"github.com/spf13/viper"

	"github.com/dgraph-io/badger/v3"
	"github.com/go-co-op/gocron"
	"github.com/streadway/amqp"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}


func GetFromRabbitMQ(db *badger.DB) {
	conn, err := amqp.Dial(viper.GetString("rabbitMQ.url"))
	if err != nil {
		fmt.Println("Failed Initializing Broker Connection")
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer ch.Close()

	if err != nil {
		fmt.Println(err)
	}

	msgs, err := ch.Consume(
		"TestQueue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	wb := db.NewWriteBatch()
	defer wb.Cancel()
	
	key := func(i int) []byte {
		return []byte(fmt.Sprintf("%d", i))
	}
	refID := 1
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Println("recieved log from RabbitMQ")
			var stud Models.Syslog
			err = json.Unmarshal(d.Body, &stud)
			stud.ID = refID
			refID = refID +1
			database.WriteToBadger(db, key(stud.ID), []byte(d.Body))
		}
	}()
	check(err)
	fmt.Println("Successfully Connected to our RabbitMQ Instance")
	fmt.Println(" [*] - Waiting for messages")
	<-forever
}


func WriteToRabbitMQ() {

	s := gocron.NewScheduler(time.UTC)

	fmt.Println("Go RabbitMQ Tutorial")
	conn, err := amqp.Dial(viper.GetString("rabbitMQ.url"))
	if err != nil {
		fmt.Println(err)
		panic(1)
	}
	defer conn.Close()

	fmt.Println("Successfully Connected to our RabbitMQ Instance")

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"TestQueue",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil  {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(q)

	go API.GetSyslog(ch)

	s.StartAsync()
	s.StartBlocking()
}