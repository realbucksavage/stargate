package stargate

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func namedWebsocketServer(name string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		client, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			log.Printf("%s: upgrade failed: %v", name, err)
			return
		}
		defer client.Close()

		err = client.WriteMessage(websocket.TextMessage, []byte(name))
		if err != nil {
			log.Printf("%s: message preparation error: %v", name, err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		err = client.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Printf("cannot send close message to client websocket: %v", err)
		}
	}
}

func TestRoundRobinWebSockets(t *testing.T) {

	maxServers := 3
	backends := make([]*httptest.Server, 0)

	for i := 1; i <= maxServers; i++ {
		n := fmt.Sprintf("server_%d", i)
		backends = append(backends, httptest.NewServer(namedWebsocketServer(n)))
	}

	defer func() {
		for _, sv := range backends {
			sv.Close()
		}
	}()

	ls := newLister(backends, "ws")
	sg, err := NewRouter(ls)
	if err != nil {
		t.Fatalf("cannot create stargate proxy: %v", err)
	}

	server := httptest.NewServer(sg)
	defer server.Close()

	for i, j := 1, 1; i < 10; i++ {
		client, _, err := websocket.DefaultDialer.Dial(toUrl(server, "ws")+"/test", nil)
		if err != nil {
			t.Fatalf("cannot create dialer: %v", err)
		}

		_, msg, err := client.ReadMessage()
		if err != nil {
			t.Fatalf("cannot read message: %v", err)
		}

		t.Logf("> iteration %d: message: %s", i, (msg))

		j++
		if j > maxServers {
			j = 1
		}
	}
}
