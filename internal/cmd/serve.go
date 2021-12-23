package cmd

import (
	"github.com/Academics-and-Career-Council/Stargazer.git/internal/services"
	"github.com/Academics-and-Career-Council/Stargazer.git/internal/database"
	"github.com/Academics-and-Career-Council/Stargazer.git/internal/api"
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
		db := database.OpenBadgerDB()
		defer db.Close()
		go database.BulkWrite(db)
		API.GetSyslog(db)
		Services.GetFromRabbitMQ(db)//for kratos
		return nil
	},
}
