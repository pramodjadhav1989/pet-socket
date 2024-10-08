package models

import "time"

type Response struct {
	StatusCode        int         `json:"status_code"`
	StatusDescription string      `json:"message_code"`
	Description       string      `json:"message"`
	Response          interface{} `json:"data"`
}

type Pong struct {
	DT time.Time `json:"time"`
}

type InputMap map[string]interface{}
