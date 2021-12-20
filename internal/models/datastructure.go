package Models

import (

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
	ID			int			`json:"id_for_ref"`
	ServiceName string		`json:"service_name"`
	StatusCode	int			`json:"status_code"`
	Severity	string		`json:"severity"`
	MsgName		string		`json:"msg_name"`
	Msg			string		`json:"msg"`
	InvokedBy	string		`json:"invoked_by"`
	Result		string		`json:"result"`
	Batch      	int    		`json:"batch"`
	Timestamp	int64		`json:"timestamp"`
}
