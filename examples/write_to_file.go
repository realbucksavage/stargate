package examples

import (
	"log"
	"net/http"
	"os"

	"github.com/realbucksavage/stargate"
	"github.com/realbucksavage/stargate/pkg/listers"
	"github.com/realbucksavage/stargate/pkg/middleware"
)

func main() {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	l := listers.Static{
		Routes: map[string][]string{
			"/ds_1": {"http://app1-sv1:8080", "http://app1-sv2:8080"},
			"/ds_2": {"http://app2-sv1:8080"},
		},
	}

	sg, err := stargate.NewProxy(l, stargate.RoundRobin, middleware.LoggerWithOutput(f))
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", &sg))
}
