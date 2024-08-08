package tasking

import (
	"context"
	"sync"
)

type Status struct {
	Cxt    context.Context
	Cancel context.CancelFunc
}
type TaskMap struct {
	sync.Mutex
	Worker map[int]Status
}

func NewTaskMap() *TaskMap {
	return &TaskMap{
		Worker: make(map[int]Status),
	}
}

func (tm *TaskMap) AddTask(id int, ctx context.Context, cancelFunc context.CancelFunc) {
	tm.Lock()
	defer tm.Unlock()
	tm.Worker[id] = Status{
		Cxt:    ctx,
		Cancel: cancelFunc,
	}
}

func (tm *TaskMap) RemoveTask(id int) {
	tm.Lock()
	defer tm.Unlock()
	delete(tm.Worker, id)
}

func (tm *TaskMap) GetTask(id int) (Status, bool) {
	tm.Lock()
	defer tm.Unlock()
	s, exists := tm.Worker[id]
	return s, exists
}
