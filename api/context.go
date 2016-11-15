package api

import (
	"encoding/json"
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/tinymailer/mailer/db"
)

const (
	contentType   = "Content-Type"
	contentBinary = "application/octet-stream"
	contentText   = "text/plain"
	contentJSON   = "application/json"
	contentHTML   = "text/html"

	defaultCharset = "UTF-8"
)

// context for request scope
type context struct {
	req    *http.Request
	res    http.ResponseWriter
	db     *mgo.Database
	params map[string]string
}

// TODO: optmize this by sync.Pool for friendly GC
// newContext create a new instance of request scope context
func newContext(r *http.Request, w http.ResponseWriter, sess *mgo.Session) *context {
	var (
		dbName = db.DB().DBName()
		params = make(map[string]string)
	)

	r.ParseForm()
	for k, v := range r.Form {
		params[k] = v[0]
	}

	return &context{
		req:    r,
		res:    w,
		db:     sess.DB(dbName),
		params: params,
	}
}

func (ctx *context) JSON(code int, data interface{}) {
	ctx.res.Header().Set(contentType, contentJSON+"; charset=UTF-8")
	ctx.res.WriteHeader(code)
	if b, err := json.Marshal(data); err != nil {
		http.Error(ctx.res, err.Error(), 500)
	} else {
		ctx.res.Write(b)
	}
}

func (ctx *context) Data(code int, data []byte) {
	ctx.res.Header().Set(contentType, contentBinary)
	ctx.res.WriteHeader(code)
	ctx.res.Write(data)
}

func (ctx *context) Redirect(url string, codes ...int) {
	code := http.StatusFound
	if len(codes) > 0 {
		code = codes[0]
	}
	http.Redirect(ctx.res, ctx.req, url, code)
}

func (ctx *context) Status(code int) {
	ctx.res.WriteHeader(code)
}

func (ctx *context) ErrNotFound(data interface{}) {
	ctx.Error(http.StatusNotFound, data)
}

func (ctx *context) ErrConflict(data interface{}) {
	ctx.Error(http.StatusConflict, data)
}

func (ctx *context) ErrBadRequest(data interface{}) {
	ctx.Error(http.StatusBadRequest, data)
}

func (ctx *context) Error(code int, data interface{}) {
	switch v := data.(type) {

	case error:
		ctx.res.WriteHeader(code)
		json.NewEncoder(ctx.res).Encode(map[string]string{
			"message": v.Error(),
		})
		return

	default:
		ctx.res.WriteHeader(code)
		json.NewEncoder(ctx.res).Encode(data)
	}
}
