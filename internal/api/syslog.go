package API

import (
	"encoding/json"
	"fmt"
	"time"

	Models "github.com/Academics-and-Career-Council/Stargazer.git/internal/models"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"gopkg.in/mcuadros/go-syslog.v2"
)

func GetSyslog(ch *amqp.Channel) { 
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.RFC5424)
	server.SetHandler(handler)
	server.ListenUDP(viper.GetString("syslog.port"))
	server.ListenTCP(viper.GetString("syslog.port"))
	server.Boot()
	//countID := 0
	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			msgs := logParts["message"].(string)
			bytVal := []byte(msgs)
			var msg interface{}
			err := json.Unmarshal(bytVal, &msg)
			if err != nil {
			 	panic(err)
			}
			temp := msg.(map[string]interface{})
			//countID = countID + 1
			Ctime := temp["time"].(string)
			layout := "2006-01-02T15:04:05Z"
			T, err := time.Parse(layout, Ctime)
			if err!=nil {
				panic(err)
			}
			var singleLog Models.Syslog
			//singleLog.ID = countID + 1
			singleLog.ServiceName = logParts["app_name"].(string)
			singleLog.Severity = temp["level"].(string)
			singleLog.Msg = temp["msg"].(string)
			singleLog.InvokedBy = singleLog.ServiceName+"@iitk.ac.in"
			singleLog.MsgName = singleLog.ServiceName+" "+singleLog.Severity
			singleLog.Result = "NA"
			singleLog.StatusCode = 500
			singleLog.Timestamp = T
			singleLog.CreatedAt = time.Now()


			data, _ := json.Marshal(singleLog)
			BodyJson := string(data)
			err = ch.Publish(
				"",
				"TestQueue",
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body: []byte(BodyJson),
				},
			)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
		}
	}(channel)

	server.Wait()
}
