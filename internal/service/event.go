package service

type EventType byte

const (
	_                     = iota //типо для созадния линейно зависимых перечислений
	EventDelete EventType = iota
	EventPut
)

type Event struct {
	Sequence  uint64
	EventType EventType
	Key       string
	Value     string
}
