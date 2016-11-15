package api

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/tinymailer/mailer/lib"
)

// GET /mail?id=xxx
func getMail(ctx *context) {
	if ctx.params["id"] == "" {
		ss, err := lib.ListMail()
		if err != nil {
			ctx.Error(500, err)
			return
		}
		ctx.JSON(200, ss)
		return
	}
}

// DELETE /mail?id=xxx
func rmMail(ctx *context) {
	id := ctx.params["id"]
	if !bson.IsObjectIdHex(id) {
		ctx.ErrBadRequest(id)
		return
	}
	bid := bson.ObjectIdHex(id)
	if err := lib.DelMail(bid); err != nil {
		ctx.Error(500, err)
		return
	}
	ctx.Status(204)
}
