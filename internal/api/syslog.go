package API

import (
	"encoding/json"
	"fmt"

	"gopkg.in/mcuadros/go-syslog.v2"
	//"github.com/Academics-and-Career-Council/Stargazer.git/internal/structure"
)

func GetSyslog() { //GetSyslog
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.RFC5424)
	server.SetHandler(handler)
	server.ListenUDP("0.0.0.0:5140")
	server.ListenTCP("0.0.0.0:5140")
	server.Boot()
	//fmt.Println("hi")

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			// syslogData, _ := json.Marshal(logParts)
			// var temp interface{}
			// err := json.Unmarshal(syslogData, &temp)
			// if err != nil {
			// 	panic(err)
			// }
			fmt.Println(logParts)
			msgs := logParts["message"].(string)
			//fmt.Println(msgs)
			bytVal := []byte(msgs)
			var msg interface{}
			err := json.Unmarshal(bytVal, &msg)
			if err != nil {
			 	panic(err)
			}
			fmt.Println(msg)
			//unmarshedSysData := string(syslogData)

			//var unmarshedSysData interface{}
			//err := json.Unmarshal(syslogData, &unmarshedSysData)
			//
			//if err != nil {
			//	panic(err)
			//}
			//fmt.Println(unmarshedSysData)
			//fmt.Println("hi")
		}
	}(channel)

	server.Wait()
}
