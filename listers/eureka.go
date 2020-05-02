package listers

import (
	"fmt"
	"github.com/hudl/fargo"
	log "github.com/op/go-logging"
	"github.com/realbucksavage/stargate"
	"strings"
)

type eurekaLister struct {
	conn fargo.EurekaConnection

	routes map[string][]string
}

func (e eurekaLister) List(s string) []string {
	return e.routes[s]
}

func (e eurekaLister) ListAll() (map[string][]string, error) {
	apps, err := e.conn.GetApps()
	if err != nil {
		return nil, err
	}

	routes := make(map[string][]string, 1)
	for _, a := range apps {
		svc := make([]string, 0)
		r := fmt.Sprintf("/%s", strings.ToLower(a.Name))
		for _, instance := range a.Instances {

			url := fmt.Sprintf("http://%s:%d", instance.IPAddr, instance.Port)
			if instance.SecurePortEnabled {
				url = fmt.Sprintf("https://%s:%d", instance.IPAddr, instance.SecurePort)
			}
			svc = append(svc, url)
		}
		routes[r] = svc
	}
	e.routes = routes

	return e.routes, nil
}

func Eureka(address string) stargate.ServiceLister {
	// Get rid of those annoying messages
	log.SetLevel(log.ERROR, "fargo")

	return eurekaLister{
		conn:   fargo.NewConn(address),
		routes: make(map[string][]string, 0),
	}
}
