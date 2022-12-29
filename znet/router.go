package znet

import "github.com/wuyutaott/b45/ziface"

type BaseRouter struct {
}

func (br *BaseRouter) PreHandle(ziface.IRequest) {}

func (br *BaseRouter) Handle(ziface.IRequest) {}

func (br *BaseRouter) PostHandle(ziface.IRequest) {}
