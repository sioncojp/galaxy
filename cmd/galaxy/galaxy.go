package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/jessevdk/go-flags"
	g "github.com/sioncojp/galaxy"
)

const (
	app = "galaxy"
)

// CommandOpts ...Load option.
type CommandOpts struct {
	Config string `long:"config" short:"c" description:"load yaml config path" required:"true"`
}

func main() {
	log.SetOutput(os.Stderr)
	log.SetPrefix(app + ": ")
	runtime.GOMAXPROCS(runtime.NumCPU())

	opts := CommandOpts{}
	if _, err := flags.ParseArgs(&opts, os.Args[1:]); err != nil {
		log.Fatalf("Error: %s", err)
	}

	config, err := g.LoadConfig(opts.Config)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	router := g.Router(config)
	router.Run(fmt.Sprintf(":%d", config.Server.Port))
}
