package proxy

import (
	"bytes"
	"github.com/nats-io/nats.go"
	go_pdk "github.com/reijiokito/go-pdk"
	"google.golang.org/protobuf/proto"
	"net/http"
	"time"
)

type Proxy struct {
	Contexts   go_pdk.HttpContextPool
	Connection *nats.Conn
}

func NewProxy(pdk *go_pdk.PDK) *Proxy {
	ret := &Proxy{
		Contexts:   go_pdk.NewContextPool(),
		Connection: pdk.Connection,
	}
	return ret
}

func (proxy *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	url := req.URL.String()

	/*CORS*/
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	ctx := proxy.Contexts.GetContext()
	defer proxy.Contexts.PutContext(ctx)

	ctx.Request.JSON = true

	/*build request*/
	err := ctx.Request.Build(req)
	if err != nil {
		http.Error(rw, "Cannot process request", http.StatusInternalServerError)
		return
	}
	/*serialize the request */
	reqBytes, err := proto.Marshal(&ctx.Request)
	if err != nil {
		http.Error(rw, "Cannot process request", http.StatusInternalServerError)
		return
	}

	/*post request to message queue*/
	msg, respErr := proxy.Connection.Request(go_pdk.PublishURL(url), reqBytes, 10*time.Second)
	if respErr != nil {
		http.Error(rw, "No response ", http.StatusInternalServerError)
		return
	}

	/*response*/
	if err := proto.Unmarshal(msg.Data, &ctx.Response); err != nil {
		http.Error(rw, "Cannot deserialize response", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(int(ctx.Response.Code))
	_, err = bytes.NewBuffer(ctx.Response.Body).WriteTo(rw)
	if err != nil {
		http.Error(rw, "Cannot deserialize response", http.StatusInternalServerError)
		return
	}
	if f, ok := rw.(http.Flusher); ok {
		f.Flush()
	}

	if f, ok := rw.(http.Flusher); ok {
		f.Flush()
	}
}
