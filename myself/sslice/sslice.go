package sslice

type SafeSlice interface {
	Append(interface{})
	At(int) interface{}
	Close() []interface{}
	Delete(int)
	Len() int
	Update(int, UpdateFunc)
}

type UpdateFunc func(interface{}) interface{}

type safeSlice chan commandData

type commandData struct {
	action	commandAction
	index	int
	value	interface{}
	result	chan<- interface{}
	data	chan<- []interface{}
	updater	UpdateFunc
}

type commandAction int

const (
	remove commandAction = iota
	end
	add
	at
	length
	update
)

func (sl safeSlice) Append(value interface{}) {
	sl <- commandData{action: add, value: value}
}

func (sl safeSlice) At(index int) interface{} {
	reply := make(chan interface{})
	sl <- commandData{action: at, index: index, result: reply}
	return <-reply
}

func (sl safeSlice) Close() []interface{} {
	reply := make(chan []interface{})
	sl <- commandData{action: end, data: reply}
	return <-reply
}

func (sl safeSlice) Delete(index int) {
	sl <- commandData{action: remove, index: index}
}

func (sl safeSlice) Len() int {
	reply := make(chan interface{})
	sl <- commandData{action: length, result: reply}
	return (<-reply).(int)
}

func (sl safeSlice) Update(index int, updater UpdateFunc) {
	sl <- commandData{action: update, index: index, updater: updater}
}

func New() SafeSlice {
	sl := make(safeSlice)
	go sl.run()
	return sl
}

func (sl safeSlice) run() {
	store := make([]interface{}, 0)
	for command := range sl {
		switch command.action {
		case add:
			store = append(store, command.value)
		case remove:
			if 0 <= command.index && command.index < len(store) {
				store = append(store[:command.index], store[command.index + 1:]...)
			}
		case at:
			if 0 <= command.index && command.index < len(store) {
				command.result <- store[command.index]
			} else {
				command.result <- nil
			}
		case length:
			command.result <- len(store)
		case update:
			if 0 <= command.index && command.index < len(store) {
				store[command.index] = command.updater(store[command.index])
			}
		case end:
			close(sl)
			command.data <- store
		}
	}
}
