package reloading_router

import (
	"log"
	"net/http"
	"time"

	"github.com/realbucksavage/stargate"
)

func main() {

	l := customLister{
		Routes: map[string][]string{
			"/test": {"http://localhost:8081"},
			"/ds_1": {"http://localhost:8082"},
		},
	}
	sg, err := stargate.NewRouter(&l)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		time.Sleep(time.Second * 25)
		l.changeRoutes()

		sg.Reload()
	}()

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":7000", sg))
}

type customLister struct {
	Routes map[string][]string
}

func (c *customLister) List(route string) ([]string, error) {
	return c.Routes[route], nil
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
