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
	fmt.Println("hi")

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			syslogData, _ := json.Marshal(logParts)
			//var unmarshedSysData interface{}
			//err := json.Unmarshal(syslogData, &unmarshedSysData)
			unmarshedSysData := string(syslogData)
			//if err != nil {
			//	panic(err)
			//}
			fmt.Println(unmarshedSysData)
			fmt.Println("hi")
		}
	}(channel)

	server.Wait()
}
