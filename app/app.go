package app

import (
	"goProj/server"
)

type app struct  {
	Server *server.Server
}

func InitApp() (*app, error) {
	serv, err := server.InitServer()
	if err != nil {
		return nil, err
	}
	return &app{
		Server: serv,
	}, nil

}

func (a *app) Run() {
	if a.Server != nil {
		a.Server.Run()
	}
}
