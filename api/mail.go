package api

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/tinymailer/mailer/lib"
)

// GET /mail?id=xxx
func getMail(ctx *context) {
	id := ctx.params["id"]
	if !bson.IsObjectIdHex(id) {
		ms, err := lib.ListMail()
		if err != nil {
			ctx.Error(500, err)
			return
		}
		ctx.JSON(200, ms)
		return
	}

	bid := bson.ObjectIdHex(id)
	m, err := lib.GetMail(bid)
	if err != nil {
		ctx.Error(500, err)
		return
	}
	ctx.JSON(200, m)
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
