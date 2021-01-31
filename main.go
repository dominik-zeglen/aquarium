package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dominik-zeglen/aquarium/api"
	"github.com/dominik-zeglen/aquarium/middleware"
	"github.com/dominik-zeglen/aquarium/sim"
	"github.com/dominik-zeglen/aquarium/tracing"
	"github.com/opentracing/opentracing-go"
)

var maxCellsInSim = flag.Int(
	"mc",
	2e4,
	"Maximum cells in the whole sim",
)
var verbose = flag.Bool(
	"v",
	false,
	"Enable verbose mode",
)
var maxCellsInOrganism = flag.Int(
	"mo",
	25,
	"Maximum cells in one organism",
)
var envDivisions = flag.Int(
	"d",
	4,
	"Divide environment along and across by this number",
)
var warmupIterations = flag.Int(
	"w",
	600,
	"Do not pause sim until number of this iterations has been reached",
)
var startCells = flag.Int(
	"s",
	10,
	"Number of cells created at the start of the sim",
)
var trace = flag.Bool(
	"t",
	false,
	"Enable tracing",
)

func getConfig() sim.SimConfig {
	return sim.SimConfig{
		EnvDivisions:       *envDivisions,
		MaxCellsInOrganism: *maxCellsInOrganism,
		MaxCellsInSim:      *maxCellsInSim,
		StartCells:         *startCells,
		Verbose:            *verbose,
		WarmupIterations:   *warmupIterations,
	}
}

func checkEnvVar(key string) error {
	if os.Getenv(key) == "" {
		return fmt.Errorf("Environment variable %s not set", key)
	}

	return nil
}

func init() {
	flag.Parse()
	envVars := []string{"ALLOWED_ORIGINS", "PORT"}

	for _, envVar := range envVars {
		err := checkEnvVar(envVar)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	if os.Getenv("JAEGER_AGENT_HOST") != "" && *trace {
		tracer, closer := tracing.InitJaeger()
		opentracing.SetGlobalTracer(tracer)
		defer closer.Close()
	}

	s := sim.Sim{}
	config := getConfig()
	s.Create(config)

	var data sim.IterationData
	http.Handle("/api",
		middleware.WithTracing(
			middleware.WithSim(
				middleware.WithCors(
					strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
					api.InitAPI(&s, &data),
				),
				&s,
			),
		),
	)
	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	s.RunLoop(&data)
}
