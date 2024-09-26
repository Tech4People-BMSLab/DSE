package debug

import (
	"os"

	"github.com/davecgh/go-spew/spew"
)

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------

// ------------------------------------------------------------
// : Functions
// ------------------------------------------------------------
func Log(v... any) {
	spew.Dump(v)
}

func Println(v... any) {
	spew.Println(v)
}

func WriteToFile(v... any) {
	file, err := os.Create("temp.bin")
	if err != nil {
		spew.Fdump(os.Stderr, err)
		return
	}
	defer file.Close()

	spew.Fdump(file, v)
}

func WriteJSONToFile(v... any) {
	file, err := os.Create("temp.json")
	if err != nil {
		spew.Fdump(os.Stderr, err)
		return
	}
	defer file.Close()

	spew.Fdump(file, v)

}

func IsDebugMode() bool {
	return os.Getenv("DEBUG") == "1"
}

func IsDebugProfiler() bool {
	return os.Getenv("DEBUG_PROFILER") == "1"
}
