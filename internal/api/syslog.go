package API

import (
	"encoding/json"
	"log"

	//"fmt"
	"time"

	"github.com/Academics-and-Career-Council/Stargazer.git/internal/database"
	Models "github.com/Academics-and-Career-Council/Stargazer.git/internal/models"
	"github.com/dgraph-io/badger/v3"
	"github.com/spf13/viper"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gopkg.in/mcuadros/go-syslog.v2"
)

func GetSyslog(db *badger.DB) {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.RFC5424)
	server.SetHandler(handler)
	server.ListenUDP(viper.GetString("syslog.port"))
	server.ListenTCP(viper.GetString("syslog.port"))
	server.Boot()

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			msgs := logParts["message"].(string) //data from recieved logs from kratos is extracted
			bytVal := []byte(msgs)
			var msg interface{}
			err := json.Unmarshal(bytVal, &msg)
			if err != nil {
				panic(err)
			}
			temp := msg.(map[string]interface{})
			Ctime := temp["time"].(string)
			layout := "2006-01-02T15:04:05Z"
			T, err := time.Parse(layout, Ctime)
			if err != nil {
				panic(err)
			}
			var singleLog Models.Syslog
			singleLog.ServiceName = logParts["app_name"].(string)
			singleLog.Severity = temp["level"].(string)
			singleLog.Msg = temp["msg"].(string)
			singleLog.InvokedBy = singleLog.ServiceName + "@iitk.ac.in"
			singleLog.MsgName = singleLog.ServiceName + " " + singleLog.Severity
			singleLog.Result = "NA"
			singleLog.StatusCode = 500
			singleLog.Timestamp = T
			singleLog.CreatedAt = time.Now()

			//All data is put into a log object
			data, _ := json.Marshal(singleLog)
			id, err := gonanoid.New()
			ID := []byte(id)
			if err != nil {
				log.Println(err)
			}
			database.WriteToBadger(db, ID, data) //object is sent to badgerDB

		}
	}(channel)

	server.Wait()
}
