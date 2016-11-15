package api

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/tinymailer/mailer/lib"
)

// GET /recipient?id=xxx
func getRecipient(ctx *context) {
	if ctx.params["id"] == "" {
		rs, err := lib.ListRecipient()
		if err != nil {
			ctx.Error(500, err)
			return
		}
		ctx.JSON(200, rs)
		return
	}
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
