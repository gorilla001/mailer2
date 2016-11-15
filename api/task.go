package api

import (
	"encoding/json"

	"gopkg.in/mgo.v2/bson"

	"github.com/tinymailer/mailer/lib"
	"github.com/tinymailer/mailer/types"
)

// GET /task?id=xxx
func getTask(ctx *context) {
	if ctx.params["id"] == "" {
		ts, err := lib.ListTask()
		if err != nil {
			ctx.Error(500, err)
			return
		}
		tws := make([]types.TaskWrapper, 0, len(ts))
		for _, t := range ts {
			tws = append(tws, lib.GetTaskWrapper(t))
		}
		ctx.JSON(200, tws)
		return
	}
}

// DELETE /task?id=xxx
func rmTask(ctx *context) {
	id := ctx.params["id"]
	if !bson.IsObjectIdHex(id) {
		ctx.ErrBadRequest(id)
		return
	}
	bid := bson.ObjectIdHex(id)

	if err := lib.DelTask(bid); err != nil {
		ctx.Error(500, err)
		return
	}
	ctx.Status(204)
}

// PATCH /task/run?id=xxx
func runTask(ctx *context) {
	id := ctx.params["id"]
	if !bson.IsObjectIdHex(id) {
		ctx.ErrBadRequest(id)
		return
	}
	bid := bson.ObjectIdHex(id)

	task, err := lib.GetTask(bid)
	if err != nil {
		ctx.Error(500, err)
		return
	}

	var opts types.TaskOptions
	if err := json.NewDecoder(ctx.req.Body).Decode(&opts); err != nil {
		ctx.Error(500, err)
		return
	}

	if err := lib.RunTask(&task, &opts); err != nil {
		ctx.Error(500, err)
		return
	}
	ctx.Status(202)
}

// PATCH /task/stop?id=xxx
func stopTask(ctx *context) {
}
