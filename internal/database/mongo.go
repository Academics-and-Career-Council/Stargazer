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
	model := mongo.IndexModel{
		Keys:    bson.M{"createdAt": 1},
		Options: options.Index().SetExpireAfterSeconds(int32((14*24*time.Hour) / time.Second)),
	}
	ind, err := database.Collection("ug").Indexes().CreateOne(context.TODO(), model)//initialising to set a TTL
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ind)
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
	data := []interface{} {}
	for role := range roles {
		data = append(data, roles[role])
	}//converting into appropriate form

	collName := viper.GetString("mongo.collName")


 	_, err = m.Logs.Collection(collName).InsertMany(context.TODO(),data)//this writes to mongoDB
	if data != nil && err == nil {
		log.Println("Inserted the batch to mongoDB")
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
	fmt.Println(batchID)
 	return batchID
}
