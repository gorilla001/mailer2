package api

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/tinymailer/mailer/lib"
)

// GET /server?id=xxx
func getServer(ctx *context) {
	id := ctx.params["id"]
	if !bson.IsObjectIdHex(id) {
		ss, err := lib.ListServer()
		if err != nil {
			ctx.Error(500, err)
			return
		}
		ctx.JSON(200, ss)
		return
	}

	bid := bson.ObjectIdHex(id)
	s, err := lib.GetServer(bid)
	if err != nil {
		ctx.Error(500, err)
		return
	}
	ctx.JSON(200, s)
}

// DELETE /server?id=xxx
func rmServer(ctx *context) {
	id := ctx.params["id"]
	if !bson.IsObjectIdHex(id) {
		ctx.ErrBadRequest(id)
		return
	}
	bid := bson.ObjectIdHex(id)
	if err := lib.DelServer(bid); err != nil {
		ctx.Error(500, err)
		return
	}
	ctx.Status(204)
}
