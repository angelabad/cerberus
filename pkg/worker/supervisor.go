package worker

import (
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"devcircus.com/cerberus/pkg/config"
	"devcircus.com/cerberus/pkg/execute"
)

// Supervisor config data
type Supervisor struct {
}

var (
	workers []Worker
)

// NewSupervisor create a new supervisor worker
func NewSupervisor() *Supervisor {
	s := Supervisor{}

	return &s
}

// Run launch the worker jobs
func (s *Supervisor) Run() {

	data := config.System.Requests
	for i, requestConfig := range data {

		println("Launching worker #", i, " : ", requestConfig.RequestType, " ", requestConfig.URL)

		w := NewWorker(requestConfig)
		workers = append(workers, *w)
		w.Start()
	}

LOOP:
	for {
		// Calling Sleep method
		time.Sleep(5 * time.Second)
		select {
		case <-execute.Done:
			log.Info().Msg("Graceful termination")
			os.Exit(0)
		case <-execute.Stop:
			log.Warn().Msg("Process terminated by external signal")
			break LOOP
		case <-execute.Reload:
			log.Info().Msg("Reloading configuration")
		default:

			log.Debug().Msg("Supervisor loop signal")

		}
	}
	os.Exit(1)
}
