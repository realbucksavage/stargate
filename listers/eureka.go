package listers

import (
	"fmt"
	"strings"

	"github.com/hudl/fargo"
	log "github.com/op/go-logging"

	"github.com/realbucksavage/stargate/v1"
)

type eurekaLister struct {
	conn fargo.EurekaConnection
}

// Eureka is an implementation of `stargate.ServiceLister` that adds support to list services registered in the passed
// eureka server. The application name of each registered service is used as a route key, e.g, if an application
// "SOME-App" is registered with Eureka, it will be accessible at /some-app.
func Eureka(address string) stargate.ServiceLister {
	// Get rid of those annoying messages
	log.SetLevel(log.ERROR, "fargo")

	return &eurekaLister{
		conn: fargo.NewConn(address),
	}
}

func (e *eurekaLister) List(s string) ([]string, error) {
	s = strings.TrimPrefix(s, "/")
	app, err := e.conn.GetApp(s)
	if err != nil {
		return nil, err
	}

	return e.listInstances(app), nil
}

// ListAll queries the underlying eureka server and creates a route table based on registered applications and their
// running instances. ListAll's implementation returns all instances with status `UP`.
func (e *eurekaLister) ListAll() (map[string][]string, error) {
	apps, err := e.conn.GetApps()
	if err != nil {
		return nil, err
	}

	routes := make(map[string][]string)
	for _, a := range apps {
		r := fmt.Sprintf("/%s", strings.ToLower(a.Name))
		routes[r] = e.listInstances(a)
	}

	return routes, nil
}

func (e *eurekaLister) listInstances(app *fargo.Application) []string {
	instances := make([]string, 0)

	for _, i := range app.Instances {
		if i.Status == fargo.UP {
			url := fmt.Sprintf(`http://%s:%d`, i.IPAddr, i.Port)
			if i.SecurePortEnabled {
				url = fmt.Sprintf(`https://%s:%d`, i.IPAddr, i.SecurePort)
			}

			instances = append(instances, url)
		}
	}

	return instances
}

var _ stargate.ServiceLister = (*eurekaLister)(nil)
