package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	//"time"

	//"github.com/Academics-and-Career-Council/Stargazer.git/internal/badgerRabbitmq"
	"github.com/Academics-and-Career-Council/Stargazer.git/internal/models"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/x/bsonx"

	//"golang.org/x/net/internal/timeseries"
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
	// index := mongo.IndexModel{
	// 	Keys:    bsonx.Doc{{Key: "created_at", Value: bsonx.Int32(1)}},
	// 	Options: options.Index().SetExpireAfterSeconds(int32(time.Now().Add(time.Second * 90).Unix())), // Will be removed after 24 Hours.
	// }
	// _, err = database.Collection("ug").Indexes().CreateOne(context.Background(), index )
	// if err != nil {
	// 	panic(err)
	// } 
	
	// database.CreateCollection(context.TODO(), "ug", {
	// 	timeseries: {
	// 		timeField: "timestamp",
	// 		metaField: "metadata",
	// 		granularity: "hours"
	// 	}
    // })
	return database
}

func (m mongoClient) BulkWriteInStudents(roles []Models.Student, db *badger.DB, bID int) error {
	var bdoc []interface{}
	docs, err := json.Marshal(roles)
	if err != nil {
    	panic(err)
	}
	err = bson.UnmarshalJSON([]byte(docs),&bdoc)
	if err != nil {
    	panic(err)
	}

	

	// myoptions := options.IndexOptions
	// myoptions
	// m.Logs.Collection("ug").Indexes().ListSpecifications(context.TODO(), {"lastModifiedDate": 3600}, )
	//m.Logs.CreateCollection("ug", {timeseries:{timeField: "timestamp"}})
	m.Logs.Collection("ug").InsertMany(context.TODO(),bdoc)
	fmt.Println("now check")
	if err != nil {
		log.Printf("Unable to check access : %v", err)
	}
	DeleteFromBadger(db,bID)
	return err
}

func (m mongoClient) GetLastBatchID() int {
	var JSONData Models.Student
	myOptions := options.FindOne()
	myOptions.SetSort(bson.M{"$natural":-1})
	lastRes := m.Logs.Collection("ug").FindOne(context.Background(), bson.M{}, myOptions)
	err := lastRes.Decode(&JSONData)
 	if err!= nil {
		//panic(err)
		return -1
	}
	//fmt.Println(JSONData)
	batchID := JSONData.Batch
 	return batchID
}
