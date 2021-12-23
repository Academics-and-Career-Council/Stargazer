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
	ServiceName string		`bson:"service_name"`
	StatusCode	int			`bson:"status_code"`
	Severity	string		`bson:"severity"`
	MsgName		string		`bson:"msg_name"`
	Msg			string		`bson:"msg"`
	InvokedBy	string		`bson:"invoked_by"`
	Result		string		`bson:"result"`
	Batch      	int    		`bson:"batch"`
	Timestamp	time.Time	`bson:"timestamp"`
	CreatedAt 	time.Time   `bson:"createdAt,omitempty"`
}
