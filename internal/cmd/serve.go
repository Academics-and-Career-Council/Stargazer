package cmd

import (
	//"fmt"
	//"log/syslog"
	"fmt"
	"time"

	"github.com/Academics-and-Career-Council/Stargazer.git/internal/services"
	"github.com/Academics-and-Career-Council/Stargazer.git/internal/database"
	"github.com/Academics-and-Career-Council/Stargazer.git/internal/models"
	"github.com/Academics-and-Career-Council/Stargazer.git/internal/api"
	"github.com/dgraph-io/badger/v3"
	"github.com/go-co-op/gocron"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the Fiber Server",
	RunE: func(cmd *cobra.Command, args []string) error {
		database.ConnectMongo()
		batchID := database.MongoClient.GetLastBatchID()
		//all related code here
		db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
		if err != nil {
			panic(err)
		}

		defer db.Close()
		//batchID := 
		//db.collection("ug").find({}, {limit: 1}).sort({$natural: -1})
		//db.Collection.find().limit(1).sort({$natural:-1})
		// wb := db.NewWriteBatch()
		// defer wb.Cancel()
		var studList []Models.Student
		s := gocron.NewScheduler(time.UTC)
		flag := false
		s.Every(30).Seconds().Do(func() {
			if !flag {
				flag = true
			} else {
				studList = database.GetFromBadger(db, batchID)
				//fmt.Println(studList)
				database.MongoClient.BulkWriteInStudents(studList, db, batchID)
				batchID = database.MongoClient.GetLastBatchID()
				fmt.Println(batchID)
			}
		})
		
		s.StartAsync()
		go API.GetSyslog()
		Services.GetFromRabbitMQ(db)
		//syslog.GetSyslog()
		s.StartBlocking()
		
		return nil
	},
}
