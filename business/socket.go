package business

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/smartpet/websocket/constant"
	"github.com/smartpet/websocket/utils"

	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/smartpet/websocket/utils/logger"
)

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	var reqID string

	reqStartTime := time.Now()
	userId := r.Header.Get(constant.USERID)

	reqID = utils.GetRequestID(r, userId)

	if utils.IsBlank(userId) {

		utils.JSONErrorResponder(r, w, http.StatusBadRequest, reqID, userId, constant.ErrorCodeMap["ABP11001"], reqStartTime, errors.New(constant.ErrorCodeMap["ABP11001"]))
		return
	}

	fmt.Println(userId)

	if !utils.ValidateJwtAndMatchClientIdCtx(r, userId) {

		w.Header().Add("Unauthorized", "true")
		utils.JSONErrorResponder(r, w, http.StatusForbidden, reqID, userId, constant.ErrorCodeMap["ABP11008"], reqStartTime, errors.New(constant.ErrorCodeMap["ABP11008"]))
		return
	}
	//validate token

	// upgrade this connection to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.ApplicationError(context.Background()).Msg(err.Error())
	}
	log.ApplicationInfo(context.Background()).Msg("Client connected")

	if err != nil {
		log.ApplicationError(context.Background()).Msg(err.Error())
	}
	reader(ws)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.ApplicationError(context.Background()).Msg(err.Error())
			return
		}
		// print out that message for clarity

		log.ApplicationInfo(context.Background()).Msg(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.ApplicationError(context.Background()).Msg(err.Error())
			return
		}
	}
}
