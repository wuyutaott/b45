package znet

import "github.com/wuyutaott/b45/ziface"

type BaseRouter struct {
}

func (br *BaseRouter) PreHandle(req ziface.IRequest) {}

func (br *BaseRouter) Handle(req ziface.IRequest) {}

func (br *BaseRouter) PostHandle(req ziface.IRequest) {}
