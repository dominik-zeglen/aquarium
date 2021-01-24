package main

import (
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

func checkEnvVar(key string) error {
	if os.Getenv(key) == "" {
		return fmt.Errorf("Environment variable %s not set", key)
	}

	return nil
}

func init() {
	envVars := []string{"ALLOWED_ORIGINS", "PORT"}

	for _, envVar := range envVars {
		err := checkEnvVar(envVar)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	if os.Getenv("JAEGER_AGENT_HOST") != "" {
		tracer, closer := tracing.InitJaeger()
		opentracing.SetGlobalTracer(tracer)
		defer closer.Close()
	}

	s := sim.Sim{}
	s.Create(os.Getenv("DEBUG") != "")

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
