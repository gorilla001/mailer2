package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/tinymailer/mailer/lib"
	"github.com/tinymailer/mailer/types"
)

// GET /task?id=xxx
func getTask(ctx *context) {
	id := ctx.params["id"]
	if !bson.IsObjectIdHex(id) {
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

	bid := bson.ObjectIdHex(id)
	t, err := lib.GetTask(bid)
	if err != nil {
		ctx.Error(500, err)
		return
	}
	ctx.JSON(200, lib.GetTaskWrapper(t))
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

// PATCH /task/run?id=xxx&sync=true
func runTask(ctx *context) {
	var (
		id = ctx.params["id"]
		//sync, _ = strconv.ParseBool(params["sync"])
	)

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

	ctx.res.Header().Set("Content-Type", "application/json")

	stream := lib.RunTask(&task, nil)
	defer stream.Close()

	var (
		msg     types.TaskProgressMsg
		decoder = json.NewDecoder(stream)
		bytes   []byte
	)
	for {
		err = decoder.Decode(&msg)
		if err == io.EOF {
			break
		}
		if err != nil {
			msg.Finish = true
			msg.Error = err.Error()
		}

		bytes, _ = json.Marshal(msg)
		bytes = append(bytes, '\r', '\n')
		ctx.res.Write(bytes)
		ctx.res.(http.Flusher).Flush()

		if msg.Finish {
			break
		}

		msg = types.TaskProgressMsg{}
	}
}

// PATCH /task/stop?id=xxx
func stopTask(ctx *context) {
	ctx.Error(500, errors.New("not implement yet"))
}
