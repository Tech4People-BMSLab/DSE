package main

import (
	"dse/src/core/log"
	"dse/src/core/services/api"
	"dse/src/core/services/api/metrics"
	"dse/src/core/services/crawler"
	"dse/src/core/services/db"
	"dse/src/core/services/extractor"
	"dse/src/core/services/monitor"
	"dse/src/core/services/scheduler"
	"dse/src/utils"
	"dse/src/utils/env"
	"dse/src/utils/event"
	"os"
	"time"

	"github.com/olebedev/emitter"
)

// TODO: Move the searches from searches_2 into the current table
// TODO: Fix issue with the pop up in the browser to start the test (trigger)
// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
type User struct {
	Firstname string
	Lastname  string
}

var (
	logger = utils.NewLogger()
)

// ------------------------------------------------------------
// : Init
// ------------------------------------------------------------
func init() {
	os.MkdirAll("tmp" , os.ModePerm)
	os.MkdirAll("logs", os.ModePerm)

	os.MkdirAll("public", os.ModePerm)
	os.MkdirAll("config", os.ModePerm)
	os.MkdirAll("data"  , os.ModePerm)
	os.MkdirAll("data/extractor", os.ModePerm)
}

// ------------------------------------------------------------
// : Debug
// ------------------------------------------------------------
func Debug() {


	// if env.IsLocal() == false { return }

	// api.Wait()

	// fsearches, _ := file.Open("searches.json")
	// defer fsearches.Close()

	// foutput, _ := file.Create("output.csv")
	// defer foutput.Close()

	// count_token  := 0 
	// count_age    := 0
	// count_sex    := 0

	// scanner := bufio.NewScanner(fsearches)
	// for scanner.Scan() {
	// 	line := scanner.Text()

	// 	timestamp := json.Get(line, "timestamp").String()
	// 	token  := json.Get(line, "token").String()
	// 	age    := json.Get(line, "form.age").String()
	// 	sex    := json.Get(line, "form.sex").String()

	// 	if token  != "" { count_token++  }
	// 	if age    != "" { count_age++ }
	// 	if sex    != "" { count_sex++ }

	// 	foutput.WriteString(fmt.Sprintf("%s, %d, %d, %d\n", timestamp, count_token, count_age, count_sex))
	// }
	// if err := scanner.Err(); err != nil {
	// 	logger.Error().Err(err).Msg("Error reading file")
	// }

	// log.Println(count_token, count_age, count_sex)
}

// ------------------------------------------------------------
// : Main
// ------------------------------------------------------------
func main() {
	log.Init() // TODO: Move this to init?
	env.Load() // TODO: Move this to init?

	log.Info().Msg("🚀 Starting...")
	
	event.Use("*", emitter.Sync, emitter.Skip)

	// Start services
	go db       .Start()
	go api      .Init()
	go crawler  .Init()
	go extractor.Init()
	go scheduler.Init()
	go monitor  .Init()
	go metrics  .Init()

	go func() {
		time.Sleep(1 * time.Hour)

		process, err := os.FindProcess(os.Getpid())
		if err != nil {
			log.Error().Err(err).Msg("Failed to find process")
			return
		}
		
		err = process.Signal(os.Kill)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send SIGKILL")
		}
	}()

	go Debug()

	select { }
}
