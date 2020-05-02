package examples

import (
	"github.com/realbucksavage/stargate"
	"github.com/realbucksavage/stargate/balancers"
	"log"
	"net/http"
	"time"
)

func main() {

	l := customLister{
		Routes: map[string][]string{
			"/test": {"http://localhost:8081"},
			"/ds_1": {"http://localhost:8082"},
		},
	}
	sg, err := stargate.NewProxy(&l, balancers.RoundRobin)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		time.Sleep(time.Second * 25)
		l.changeRoutes()

		sg.Reload()
	}()

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":7000", &sg))
}

type customLister struct {
	Routes map[string][]string
}

func (c *customLister) List(route string) []string {
	return c.Routes[route]
}

func (c *customLister) ListAll() (map[string][]string, error) {
	return c.Routes, nil
}

func (c *customLister) changeRoutes() {
	c.Routes = map[string][]string{
		"/":     {"http://localhost:8082", "http://localhost:8083"},
		"/ds_2": {"http://localhost:8081"},
	}
}
