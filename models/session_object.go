package models

import (
	"time"
)

type SessionObject struct {
	Name       string      `json:"name"`
	Value      interface{} `json:"value"`
	CreateTime time.Time   `json:"create_time"`
	UpdateTime time.Time   `json:"update_time"`
}
