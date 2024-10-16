package triple_s

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ab-dauletkhan/triple-s/api"
	"github.com/ab-dauletkhan/triple-s/api/core"
	"github.com/ab-dauletkhan/triple-s/api/util"
)

func Run() {
	// Parses the port, dir and help flags.
	// If, help provided prints help message immediately and program stops there
	core.ParseFlags()
	if core.Help {
		core.PrintUsage()
		return
	}

	err := util.InitDir()
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", core.Port),
		Handler: api.Routes(),
	}

	log.Printf("Starting the server on %d...\n", core.Port)
	log.Printf("Data dir: %s", core.Dir)
	err = srv.ListenAndServe()
	log.Fatal(err)
}
