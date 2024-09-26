package scheduler

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"bms.dse/src/common"
	"bms.dse/src/services/crawler"
	"bms.dse/src/services/ipc"
	"bms.dse/src/utils"
	"bms.dse/src/utils/logutil"
	"github.com/robfig/cron/v3"
)

// ------------------------------------------------------------
// : Types
// ------------------------------------------------------------
type Client = common.Client

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var logger = logutil.NewLogger("Scheduler")
var scheduler cron.Cron
// ------------------------------------------------------------
// : Handlers
// ------------------------------------------------------------
func Consent() {
	users := crawler.GetUsers()

	for _, user := range users {
		subject := fmt.Sprintf("bms.dse.users.%s.consent", user.Token)
		ipc.Publish(subject, []byte("{}"))
	}
}

func Trigger() {
	go func() {
		users := crawler.GetUsers()
	
		for _, user := range users {
			user.Start()
		}
	}()
}

func Restart() {
	logger.Info().Msg("Restarting...")
	os.Exit(0)
}

func GC() {
	logger.Info().Msg("GC triggered")
	runtime.GC()
}
// ------------------------------------------------------------
// : Functions	
// ------------------------------------------------------------
func AddFunc(spec string, cmd func()) {
	scheduler.AddFunc(spec, cmd)
}
// ------------------------------------------------------------
// : Scheduler
// ------------------------------------------------------------
func Init() {
	if utils.IsDebugMode() { return }

	scheduler = *cron.New(
		cron.WithLocation(time.UTC),
	)

	scheduler.AddFunc("0 4 * * *", Restart)   // Every day at 4am
	scheduler.AddFunc("0 * * * *", GC)        // Every hour
	scheduler.AddFunc("0 0 12 6 9 2024", Trigger) // 12:00:00 6th September 2024

	scheduler.Start()
}
