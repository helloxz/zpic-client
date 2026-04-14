package core

import (
	"context"
	"runtime"
)

var VERSION = "1.1.0"

type AppCore struct {
	ctx context.Context
}

func NewAppCore() *AppCore {
	return &AppCore{}
}

type appInfos struct {
	Version string `json:"version"`
	OS      string `json:"os"`
}

func (ac *AppCore) GetAppInfo() ResData {
	os := runtime.GOOS
	version := VERSION
	resData := appInfos{
		Version: version,
		OS:      os,
	}
	return ResData{
		Status: true,
		Msg:    "获取成功",
		Data:   resData,
	}
}
