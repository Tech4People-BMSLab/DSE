package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"bms.dse/src/common"
	"bms.dse/src/debug"
	"bms.dse/src/services/api"
	"bms.dse/src/services/crawler"
	"bms.dse/src/services/db"
	"bms.dse/src/services/ipc"
	"bms.dse/src/services/scheduler"

	"bms.dse/src/utils/logutil"
	"github.com/felixge/fgprof"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

// ------------------------------------------------------------
// : Types
// ------------------------------------------------------------
type Client = common.Client

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var logger = logutil.NewLogger("Main")

// ------------------------------------------------------------
// : Functiins
// ------------------------------------------------------------
func LoadEnv() {
	godotenv.Load("./.env")
	return
}

func GetBuild() string {
	b, err := os.ReadFile("./build.txt")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to read build.txt")
		return "unknown"
	}
	return string(b)
}

func Debug() {
	if !debug.IsDebugMode() { return }
	log.Println("🔧 Debug mode enabled")

	if debug.IsDebugProfiler() {
		http.DefaultServeMux.Handle("/debug/fgprof", fgprof.Handler()) // go tool pprof --http=:6061 http://localhost:6060/debug/fgprof?seconds=3
		go func() {
			log.Println(http.ListenAndServe(":6060", nil))
		}()
	}

	c := resty.New()
	c.SetDebug(true)
	c.SetTimeout(100 * time.Millisecond)
	c.R().Get("http://localhost:5000/api/system/shutdown?token=dse2023")
}

// ------------------------------------------------------------
// : Main
// ------------------------------------------------------------
func main() {
	LoadEnv()
	logger.Println("🚀 Starting DSE...")
	logger.Println("📦 Build:", GetBuild())

	Debug()

	// Start modules
	go db.Init()
	go api.Init()
	go ipc.Init()
	go scheduler.Init()
	go crawler.Init()

	// Keep alive
	select {}
}
