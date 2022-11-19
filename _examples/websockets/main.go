package main

import (
	"github.com/realbucksavage/stargate"
	"github.com/realbucksavage/stargate/listers"

	"log"
	"net/http"
)

func main() {
	l := listers.Static{
		Routes: map[string][]string{
			"/": {"ws://wsserver1:8080", "ws://wsserver2:8080"},
		},
	}
	sg, err := stargate.NewRouter(l)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8081", sg))
}
