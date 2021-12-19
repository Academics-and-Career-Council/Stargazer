package Models

import (
	"time"

	"go.mongodb.org/mongo-driver/x/bsonx"
)

type Student struct {
	ID         int    		`json:"id"`
	Name       string 		`json:"name"`
	Branch     string 		`json:"branch"`
	Age        int    		`json:"age"`
	Batch      int    		`json:"batch"`
	CreatedAt  bsonx.Val	`json:"created_at"`
}

type Syslog struct {
	Service  	string		`json:"service"`
	SeverityLvl	string		`json:"severityLvl"`
	Msg			string		`json:"msg"`
	TimeStamp 	time.Time	`json:"timestamp"`
}
