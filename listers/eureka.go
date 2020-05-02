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

// ListAll queries the underlying eureka server and creates a route table based on registered applications and their
// running instances. ListAll's implementation returns all instance regardless of their health status. Healthchecks
// are done by Stargate.
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

			// Use HTTPs path if the instance supports SSL.
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

// Eureka is an implementation of `stargate.ServiceLister` that adds support to list services registered in the passed
// eureka server. The application name of each registered service is used as a route key, e.g, if an application
// "SOME-App" is registered with Eureka, it will be accessible at /some-app.
func Eureka(address string) stargate.ServiceLister {
	// Get rid of those annoying messages
	log.SetLevel(log.ERROR, "fargo")

	return eurekaLister{
		conn:   fargo.NewConn(address),
		routes: make(map[string][]string, 0),
	}
}
