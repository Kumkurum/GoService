package service

type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key string, value string)
	Error() <-chan error
	ReadEvents() (<-chan Event, <-chan error)

	Run()
}
