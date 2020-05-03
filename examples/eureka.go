package examples

import (
	"github.com/realbucksavage/stargate"
	"log"
	"net/http"
	"time"
)

func main() {
	el := stargate.EurekaLister("http://localhost:8761/eureka")
	sg, err := stargate.NewProxy(el, stargate.RoundRobin)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			time.Sleep(time.Second * 30)
			err := sg.Reload()
			if err != nil {
				log.Fatalf("Cannot reload eureka lister : %v", err)
			}
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", &sg))
}
