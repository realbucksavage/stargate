package examples

import (
	"github.com/realbucksavage/stargate"
	"log"
	"net/http"
)

func main() {
	l := stargate.StaticLister{
		Routes: map[string][]string{
			"/ds_1": {"http://app1-sv1:8080", "http://app1-sv2:8080"},
			"/ds_2": {"http://app2-sv1:8080"},
		},
	}
	sg, err := stargate.NewProxy(l, stargate.RoundRobin)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", &sg))
}
