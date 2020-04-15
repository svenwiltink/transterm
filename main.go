package main

import (
	"flag"
	"github.com/svenwiltink/transterm/app"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"

	_ "net/http/pprof"
)

// Show a navigable tree view of the current directory.
func main() {

	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		runtime.SetBlockProfileRate(1)
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()

		go func() {
			err := http.ListenAndServe("localhost:6060", nil)
			if err != nil {
				panic(err)
			}
		}()
	}
	application := app.NewApplication(app.Config{
		AccountName:    "swiltink",
		PrivateKeyPath: "transip.key",
	})

	application.Run()
}
