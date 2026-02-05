package core

import (
	"context"
)

type AppCore struct {
	ctx context.Context
}

func NewAppCore() *AppCore {
	return &AppCore{}
}
