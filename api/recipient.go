package api

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/tinymailer/mailer/lib"
)

// GET /recipient?id=xxx
func getRecipient(ctx *context) {
	id := ctx.params["id"]
	if !bson.IsObjectIdHex(id) {
		rs, err := lib.ListRecipient()
		if err != nil {
			ctx.Error(500, err)
			return
		}
		ctx.JSON(200, rs)
		return
	}

	bid := bson.ObjectIdHex(id)
	r, err := lib.GetRecipient(bid)
	if err != nil {
		ctx.Error(500, err)
		return
	}
	ctx.JSON(200, r)
}

// DELETE /recipient?id=xxx
func rmRecipient(ctx *context) {
	id := ctx.params["id"]
	if !bson.IsObjectIdHex(id) {
		ctx.ErrBadRequest(id)
		return
	}
	bid := bson.ObjectIdHex(id)
	if err := lib.DelRecipient(bid); err != nil {
		ctx.Error(500, err)
		return
	}
	ctx.Status(204)
}
