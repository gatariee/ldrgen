package main

import (
	"context"
	ldr "ldr/cmd"

)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetConfig(cfg_path string) *ldr.Config {
	config, err := ldr.ReadConfig(cfg_path)
	if err != nil {
		panic(err)
	}
	return config
}

func (a *App) ReadFile(filepath string) (string, error) {

	base_template := "../templates/"

	content, err := ldr.ReadFile(base_template + filepath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}