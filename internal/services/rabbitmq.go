package Services

import (
	"encoding/json"
	"fmt"

	"github.com/Academics-and-Career-Council/Stargazer.git/internal/database"
	"github.com/Academics-and-Career-Council/Stargazer.git/internal/models"

	"github.com/dgraph-io/badger/v3"
	"github.com/streadway/amqp"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}


func GetFromRabbitMQ(db *badger.DB) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672")
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

	forever := make(chan bool)
	go func() {
		//go syslog.GetSyslog()
		for d := range msgs {
			//fmt.Printf("Recieved Message: %s\n", d.Body)
			var stud Models.Student
			err = json.Unmarshal(d.Body, &stud)
			//stud.Batch = bID + 1
			//temp, err := json.Marshal(stud)
			//check(err)
			database.WriteToBadger(db, key(stud.ID), []byte(d.Body))
			// db.Update(func(txn *badger.Txn) error {
			// 	err := txn.SetEntry(badger.NewEntry(key(stud.ID), []byte(d.Body)))//to add Withttl
			// 	//fmt.Println("sent to badgerDB", d.Body)
			// 	return err
			// })
		}
	}()
	check(err)
	fmt.Println("Successfully Connected to our RabbitMQ Instance")
	fmt.Println(" [*] - Waiting for messages")
	<-forever
}
