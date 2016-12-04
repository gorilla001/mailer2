package api

import (
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/tinymailer/mailer/db"
)

const (
	// APIPREFIX is exported
	APIPREFIX = "/api"
)

var (
	// global router holder: method -> path -> handleFunc
	smux map[string]map[string]handleFunc
)

type handleFunc func(*context)

// ListenAndServe start up http api
func ListenAndServe(listen string) error {
	if listen == "" {
		listen = ":80"
	}

	server := http.Server{
		Addr:    listen,
		Handler: new(handler),
	}

	setupRouters()
	return server.ListenAndServe()
}

// init and setup global router store
// TODO support subgroup-midware for specified prefix routes
func setupRouters() {
	smux = map[string]map[string]handleFunc{
		"GET": map[string]handleFunc{
			"":           listAPI,
			"/":          listAPI,
			"/ping":      ping,
			"/server":    getServer,
			"/recipient": getRecipient,
			"/mail":      getMail,
			"/task":      getTask,
		},
		"DELETE": map[string]handleFunc{
			"/server":    rmServer,
			"/recipient": rmRecipient,
			"/mail":      rmMail,
			"/task":      rmTask,
		},
		"PATCH": map[string]handleFunc{
			"/task/run":  runTask,
			"/task/stop": stopTask,
		},
	}
}

// dispath request route to right handler within global route store
// TODO support to get all of subgroup-midware handlers for specified method/path
func route(method, path string) (handleFunc, bool) {
	// must prefix with APIPREFIX
	if !strings.HasPrefix(path, APIPREFIX) {
		return nil, false
	}
	path = strings.TrimPrefix(path, APIPREFIX)
	handler, has := smux[method][path]
	return handler, has
}

type handler struct{}

func (*handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// global midware ...

	// log midware
	method, remote, path := r.Method, r.RemoteAddr, r.URL.Path
	log.Println("Request:", method, remote, path)

	// auth midware

	// panic recovery midware

	// dispatch request
	if handler, has := route(method, path); has {
		dbSess := db.DB().NewSession()
		ctx := newContext(r, w, dbSess)
		handler(ctx)
		dbSess.Close()
		return
	}

	// 404
	w.WriteHeader(404)
	w.Write([]byte("handler not found"))
}

func ping(ctx *context) {
	ctx.res.Write([]byte{'O', 'K'})
}

func listAPI(ctx *context) {
	data := make(map[string]map[string]bool)
	for m, ph := range smux {
		minner := make(map[string]bool)
		for p := range ph {
			minner[p] = true
		}
		data[m] = minner
	}
	ctx.JSON(200, data)
}
