package main

import (
	"log"
	"net/http"

	"github.com/realbucksavage/stargate"
	"github.com/realbucksavage/stargate/listers"
)

type exampleLogger struct{}

func (e exampleLogger) Info(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (e exampleLogger) Warn(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (e exampleLogger) Debug(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (e exampleLogger) Error(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func main() {
	stargate.Log = exampleLogger{}

	l := listers.Static{
		Routes: map[string][]string{
			"/ds_1": {"http://app1-sv1:8080", "http://app1-sv2:8080"},
			"/ds_2": {"http://app2-sv1:8080"},
		},
	}
	sg, err := stargate.NewRouter(l)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", sg))
}
