package core

import (
	"context"
)

var VERSION = "0.1.0"

type AppCore struct {
	ctx context.Context
}

func NewAppCore() *AppCore {
	return &AppCore{}
}
