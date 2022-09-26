package Services


import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Academics-and-Career-Council/Stargazer.git/internal/database"
	"github.com/Academics-and-Career-Council/Stargazer.git/internal/models"
	gonanoid "github.com/matoous/go-nanoid/v2"
	// "github.com/spf13/viper"

	"github.com/dgraph-io/badger/v3"
	// "github.com/streadway/amqp"
	zmq "github.com/pebbe/zmq4"
)

// func check(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }


func GetFromZeroMQ(db *badger.DB) {
	// conn, err := amqp.Dial(viper.GetString("rabbitMQ.url"))
	// if err != nil {
	// 	fmt.Println("Failed Initializing Broker Connection")
	// 	panic(err)
	// }
	// defer conn.Close()

	// ch, err := conn.Channel()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer ch.Close()

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// msgs, err := ch.Consume(
	// 	"TestQueue",
	// 	"",
	// 	true,
	// 	false,
	// 	false,
	// 	false,
	// 	nil,
	// )
	zctx, _ := zmq.NewContext()

	s, err := zctx.NewSocket(zmq.REP)
	s.Bind("tcp://*:5555")

	wb := db.NewWriteBatch()
	defer wb.Cancel()
	
	
	forever := make(chan bool)
	go func() {
		for  { // d := range msgs
			log.Println("recieved log from RabbitMQ")
			msg, _ := s.Recv(0)
			log.Println(msg)
			go s.Send("Received", 0)
			var stud Models.Syslog
			err = json.Unmarshal([]byte(msg), &stud)
			stud.Timestamp = time.Now()
			binMsg, _ := json.Marshal(stud)
			id, err := gonanoid.New()
			log.Println(id)
			ID := []byte(id)
			if err!=nil {
				log.Println(err)
			}
			database.WriteToBadger(db, ID, binMsg)
		}
	}()
	check(err)
	fmt.Println("Successfully Connected to our RabbitMQ Instance")
	fmt.Println(" [*] - Waiting for messages")
	<-forever
}


