package Services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Academics-and-Career-Council/Stargazer.git/internal/database"
	"github.com/Academics-and-Career-Council/Stargazer.git/internal/models"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/spf13/viper"

	"github.com/dgraph-io/badger/v3"
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
	

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Println("recieved log from RabbitMQ")
			var stud Models.Syslog
			err = json.Unmarshal(d.Body, &stud)
			stud.Timestamp = time.Now()
			d.Body, err = json.Marshal(stud)
			id, err := gonanoid.New()
			ID := []byte(id)
			if err!=nil {
				log.Println(err)
			}
			database.WriteToBadger(db, ID, []byte(d.Body))
		}
	}()
	check(err)
	fmt.Println("Successfully Connected to our RabbitMQ Instance")
	fmt.Println(" [*] - Waiting for messages")
	<-forever
}


