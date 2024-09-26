package livelog

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"bms.dse/src/utils/logutil"
	"github.com/golang-module/carbon"
	"github.com/hpcloud/tail"
)

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var logger = logutil.NewLogger("API")
// ------------------------------------------------------------
// : Functions
// ------------------------------------------------------------
func HandleLogDownload(w http.ResponseWriter, r *http.Request) {
	paramLimit := r.URL.Query().Get("limit")
	if paramLimit == "" { paramLimit = "100" }

	filePath        := "logs/main.log"
	fileReader, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}

	fileScanner := bufio.NewScanner(fileReader)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		_, _ = w.Write([]byte(line + "\n"))
	}
	return
}

func HandleLogStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return 
	}

	paramLimit := r.URL.Query().Get("limit")
	if paramLimit == "" { paramLimit = "100" }

	w.Header().Set("Content-Type"     , "text/event-stream")
    w.Header().Set("Cache-Control"    , "no-cache")
    w.Header().Set("Connection"       , "keep-alive")
    w.Header().Set("X-Accel-Buffering", "no") // Disable buffering for nginx

	ctx, cancel := context.WithTimeout(r.Context(), 1 * time.Hour)
	defer cancel()

	timeBegin := carbon.Now()
	timeNow   := carbon.Carbon{}

	fileTail, err := tail.TailFile("logs/main.log", tail.Config{Follow: true, ReOpen: true, Poll: true})
	if err != nil {
		http.Error(w, "Failed to tail file", http.StatusInternalServerError)
		return
	}

	go func() {
		<-ctx.Done()
		fileTail.Stop()
	}()

	filePath        := "logs/main.log"
	fileReader, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to open file")
		return
	}

	fileScanner := bufio.NewScanner(fileReader)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		fmt.Fprintf(w, "data: %s\n\n", line)
		flusher.Flush()
	}

	for {
		select {
			case fileLine := <- fileTail.Lines: {
				if fileLine.Err != nil {
					http.Error(w, "Failed to read line", http.StatusInternalServerError)
					return
				}

				object := map[string]any{}
				err    := json.Unmarshal([]byte(fileLine.Text), &object)
				if err != nil {
					continue
				}

				timeNow = carbon.Parse(object["time"].(string))
				if timeNow.Lt(timeBegin) {
					continue
				}

				fmt.Fprintf(w, "data: %s\n\n", fileLine.Text)
				flusher.Flush()
			}

			case <- ctx.Done(): {
				return
			}
		}
	}

}

// ------------------------------------------------------------
// : Handler
// ------------------------------------------------------------
func HandleLog(w http.ResponseWriter, r *http.Request) {
	paramMethod := r.URL.Query().Get("method")
	paramToken  := r.URL.Query().Get("token")

	if paramToken != "2023" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch (paramMethod) {
		case "stream": {
			HandleLogStream(w, r)
			break
		}

		case "download": fallthrough
		default: {
			HandleLogDownload(w, r)
			break
		}
	}
	return
}
