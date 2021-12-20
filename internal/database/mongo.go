package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Academics-and-Career-Council/Stargazer.git/internal/models"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/dgraph-io/badger/v3"
	"gopkg.in/mgo.v2/bson"
)



type mongoClient struct {
	Logs *mongo.Database
}

var MongoClient = &mongoClient{}

func ConnectMongo() {
	MongoClient.Logs = connect(viper.GetString("mongo.url"), viper.GetString("mongo.database"))
}

func connect(url string, dbname string) *mongo.Database {
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Unable to Connect to MongoDB %v", err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Unable to Connect to MongoDB %v", err)
	}
	log.Printf("Connected to MongoDB! URL : %s", url)
	database := client.Database(dbname)
	return database
}

func (m mongoClient) BulkWriteInSyslog(roles []Models.Syslog, db *badger.DB, bID int) error {
	var bdoc []interface{}
	docs, err := json.Marshal(roles)
	if err != nil {
    	panic(err)
	}
	err = bson.UnmarshalJSON([]byte(docs),&bdoc)
	if err != nil {
    	panic(err)
	}
	index := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "created_at", Value: bsonx.Int32(1)}},
		Options: options.Index().SetExpireAfterSeconds(int32(time.Duration(90))), 
	}
	collName := viper.GetString("mongo.collName")
	_, err = m.Logs.Collection(collName).Indexes().CreateOne(context.Background(), index )
	if err != nil {
		panic(err)
	}

	m.Logs.Collection(collName).InsertMany(context.TODO(),bdoc)
	fmt.Println("Successfully Sent a batch to mongoDB")
	if err != nil {
		log.Printf("Unable to check access : %v", err)
	}
	DeleteFromBadger(db,bID)
	return err
}

func (m mongoClient) GetLastBatchID() int {
	collName := viper.GetString("mongo.collName")
	var JSONData Models.Syslog
	myOptions := options.FindOne()
	myOptions.SetSort(bson.M{"$natural":-1})
	lastRes := m.Logs.Collection(collName).FindOne(context.Background(), bson.M{}, myOptions)
	err := lastRes.Decode(&JSONData)
 	if err!= nil {
		return -1
	}
	batchID := JSONData.Batch
 	return batchID
}
