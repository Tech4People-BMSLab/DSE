package profiling

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"runtime/pprof"

	"github.com/rs/zerolog/log"
)

// ------------------------------------------------------------
// : Functions
// ------------------------------------------------------------
func Start() {
	file, err := os.Create("./cpu.pprof")
	if err != nil {
		log.Error().Err(err).Msg("Failed to create profiling file")
		return
	}
	pprof.StartCPUProfile(file)
}

func Stop() {
	pprof.StopCPUProfile()

	go func() {
		cmd := exec.Command("go", "tool", "pprof", "-http=10.0.0.10:8080", "./cpu.pprof")
		err := cmd.Run()
		if err != nil {
			log.Error().Err(err).Msg("Failed to run go tool pprof")
		}
	}()
}

// ------------------------------------------------------------
// : Initialization
// ------------------------------------------------------------
func Init() {
	go func ()  {
		err := http.ListenAndServe(":6060", nil)
		if err != nil {
			log.Error().Err(err).Msg("Failed to start profiling server")
			return
		}
	}()
}
