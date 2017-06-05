// buntobot is the server which controls @buntobot on GitHub.
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/buntobot/auto-reply/ctx"
	"github.com/buntobot/auto-reply/bunto"
)

var context *ctx.Context

func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "The port to serve to")
	flag.Parse()
	context = ctx.NewDefaultContext()

	http.HandleFunc("/_ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("ok\n"))
	}))

	buntoOrgHandler := bunto.NewBuntoOrgHandler(context)
	http.Handle("/_github/bunto", buntoOrgHandler)

	log.Printf("Listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
