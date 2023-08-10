package stargate

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type websocketDownstream struct {
	url      string
	director func(*http.Request)
}

func (w *websocketDownstream) Healthy(_ context.Context) error {
	// TODO: Implement websocket healthchecks
	return nil
}

func (w *websocketDownstream) Address() string {
	return w.url
}

func (w *websocketDownstream) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	requestHeader := http.Header{}
	if origin := req.Header.Get("Origin"); origin != "" {
		requestHeader.Add("Origin", origin)
	}

	for _, prot := range req.Header[http.CanonicalHeaderKey("Sec-WebSocket-Protocol")] {
		requestHeader.Add("Sec-WebSocket-Protocol", prot)
	}

	for _, cookie := range req.Header[http.CanonicalHeaderKey("Cookie")] {
		requestHeader.Add("Cookie", cookie)
	}

	if req.Host != "" {
		requestHeader.Set("Host", req.Host)
	}

	// Pass X-Forwarded-For headers too, code below is a part of
	// httputil.ReverseProxy. See http://en.wikipedia.org/wiki/X-Forwarded-For
	// for more information
	// TODO: use RFC7239 http://tools.ietf.org/html/rfc7239
	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		// If we aren't the first proxy retain prior
		// X-Forwarded-For information as a comma+space
		// separated list and fold multiple headers into one.
		if prior, ok := req.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		requestHeader.Set("X-Forwarded-For", clientIP)
	}

	// Set the originating protocol of the incoming HTTP request. The SSL might
	// be terminated on our site and because we doing proxy adding this would
	// be helpful for applications on the backend.
	requestHeader.Set("X-Forwarded-Proto", "http")
	if req.TLS != nil {
		requestHeader.Set("X-Forwarded-Proto", "https")
	}

	cloned := req.Clone(req.Context())
	w.director(cloned)

	destAddress := w.url
	if cloned.URL.Path != "" {
		destAddress = destAddress + cloned.URL.Path
	}

	downstreamConnection, downstreamResp, err := websocket.DefaultDialer.Dial(destAddress, requestHeader)
	if err != nil {
		Log.Debug("cannot connect to downstream server: %v", err)
		if downstreamResp != nil {
			if err := copyResponse(rw, downstreamResp); err != nil {
				Log.Error("cannot write response after handshake failure: %v", err)
			}
		} else {
			http.Error(rw, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
	}
	defer func() {
		if err := downstreamConnection.Close(); err != nil {
			Log.Debug("cannot close connection to %q: %v", downstreamConnection.RemoteAddr(), err)
		}
	}()

	upgradeHeader := http.Header{}
	if hdr := downstreamResp.Header.Get("Sec-Websocket-Protocol"); hdr != "" {
		upgradeHeader.Set("Sec-Websocket-Protocol", hdr)
	}
	if hdr := downstreamResp.Header.Get("Set-Cookie"); hdr != "" {
		upgradeHeader.Set("Set-Cookie", hdr)
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		EnableCompression: true,
	}

	clientConn, err := upgrader.Upgrade(rw, req, upgradeHeader)
	if err != nil {
		Log.Error("cannot upgrade: %v", err)
		http.Error(rw, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}

	defer func() {
		if err := clientConn.Close(); err != nil {
			Log.Error("cannot close connection to %q: %v", clientConn.RemoteAddr(), err)
		}
	}()

	var (
		downstreamErrors = make(chan error, 1)
		clientErrors     = make(chan error, 1)
	)

	go replicateWebsocketConnection(clientConn, downstreamConnection, clientErrors, "client to downstream")
	go replicateWebsocketConnection(downstreamConnection, clientConn, downstreamErrors, "downstream to client")

	var message string
	select {
	case err = <-downstreamErrors:
		message = "error copying from downstream to client: %v"
	case err = <-clientErrors:
		message = "error copying from client to downstream: %v"
	}

	if e, ok := err.(*websocket.CloseError); !ok || e.Code == websocket.CloseAbnormalClosure {
		Log.Error(message, err)
	}
}

func replicateWebsocketConnection(dest, src *websocket.Conn, errs chan error, alias string) {

	for {
		msgType, msg, err := src.ReadMessage()
		if err != nil {
			m := websocket.FormatCloseMessage(websocket.CloseNormalClosure, fmt.Sprintf("%v", err))
			if e, ok := err.(*websocket.CloseError); ok {
				if e.Code != websocket.CloseNoStatusReceived {
					m = websocket.FormatCloseMessage(e.Code, e.Text)
				}
			}

			errs <- err
			if err := dest.WriteMessage(websocket.CloseMessage, m); err != nil {
				Log.Error("error write error:: %s:: %v", alias, err)
			}

			break
		}

		err = dest.WriteMessage(msgType, msg)
		if err != nil {
			Log.Error("write error:: %s:: %v", alias, err)
			errs <- err
			break
		}
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func copyResponse(rw http.ResponseWriter, resp *http.Response) error {
	copyHeader(rw.Header(), resp.Header)
	rw.WriteHeader(resp.StatusCode)
	defer resp.Body.Close()

	_, err := io.Copy(rw, resp.Body)
	return err
}
