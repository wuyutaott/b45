package znet

import (
	"fmt"
	"github.com/wuyutaott/b45/utils"
	"github.com/wuyutaott/b45/ziface"
)

type MsgHandler struct {
	Apis           map[uint32]ziface.IRouter
	WorkerPoolSize uint32
	TaskQueue      []chan ziface.IRequest
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GCfg.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GCfg.WorkerPoolSize),
	}
}

func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	mh.TaskQueue[workerID] <- request
}

func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Printf("api msgID = %d not found!", request.GetMsgID())
		return
	}
	request.BindRouter(handler)
	request.Next()
}

func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		panic(fmt.Sprintf("repeat api err! msgID = %d \n", msgID))
	}
	mh.Apis[msgID] = router
	fmt.Println("Add api success! msgID =", msgID)
}

func (mh *MsgHandler) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Printf("Worker id = %d is start! \n", workerID)
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GCfg.MaxWorkerTaskLen)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}
