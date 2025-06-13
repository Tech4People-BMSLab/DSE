package scheduler

import (
	"dse/src/core/services/api/download"
	"dse/src/core/services/crawler"
	"dse/src/core/services/db"
	"dse/src/core/services/monitor"
	"dse/src/utils"
	"encoding/json"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var (
	logger = utils.NewLogger()
)

// ------------------------------------------------------------
// : Tasks
// ------------------------------------------------------------
func reset() {
	logger.Info().Msg("Resetting all tasks")
	crawler.OnResetAll()
}

func Consent() {
	users, err := db.GetUsers()
	if err != nil {
		panic(err)
	}

	// Prepare data to write to file
	user_data := []map[string]interface{}{}
	users.Each(func(i int, v *db.User) bool {
		user_data = append(user_data, map[string]interface{}{
			"token":  v.Token, // Changed "id" to "token"
			"online": v.IsOnline(),
		})

		if v.IsOnline() { // Check if user is online
			if v.State.Client.User.Form.Postcode.Value == "" {
				v.Send("consent", nil) // Send "consent" to the user with correct arguments
			}
		}
		return true
	})

	// Write users to /tmp/YYYYMMDD_HHMMSS_users.json
	timestamp := time.Now().Format("20060102_150405")
	file_path := "tmp/" + timestamp + "_users.json"

	file, err := os.Create(file_path) // Use os.Create to create the file
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create file")
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file) // Use json.NewEncoder to write JSON
	err = encoder.Encode(user_data)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to write users to file")
	} else {
		logger.Info().Msgf("Users written to %s", file_path)
	}
}

func debug() {
	
}


// ------------------------------------------------------------
// : Scheduler
// -----------------------------------------------------------
func Init() {
	// Triggers that should happen on start
	go download.LoadData()

	// Triggers that should happen on a schedule
	c := cron.New()
	c.AddFunc("0 0 * * 1", func() { reset() })
	c.AddFunc("*/10 * * * *", func() { monitor.MonitorUsers() })
	c.AddFunc("*/10 * * * *", func() { monitor.MonitorSearches() })
	c.AddFunc("*/10 * * * *", func() { monitor.MonitorSearchesSize() })
	c.AddFunc("*/10 * * * *", func() { monitor.MonitorSearchesTotal() })
	c.AddFunc("*/10 * * * *", func() { download.LoadData() })
	c.AddFunc("0 12 20 3 *" , func() { Consent() }) // On March 20th at 12:00 PM
	c.AddFunc("0 12 21 3 *" , func() { Consent() }) // On March 21st at 12:00 PM
	c.Start()

	go debug()
}
