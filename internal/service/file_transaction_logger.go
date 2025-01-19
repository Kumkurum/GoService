package service

import (
	"bufio"
	"fmt"
	"os"
)

type FileTransactionLogger struct {
	events       chan<- Event // Канал только для записи; для передачи событий
	errors       <-chan error // Канал только для чтения; для приема ошибок
	lastSequence uint64       // Последний использованный порядковый номер
	file         *os.File     // Местоположение файла журнала
}

func NewFileTransactionLogger(filename string) (TransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return nil, fmt.Errorf("error opening file %w", err)
	}
	return &FileTransactionLogger{file: file}, nil
}

func (f *FileTransactionLogger) WritePut(key, value string) {
	f.events <- Event{EventType: EventPut, Key: key, Value: value}
}
func (f *FileTransactionLogger) WriteDelete(key string) {
	f.events <- Event{EventType: EventDelete, Key: key, Value: "deleted"}
}
func (f *FileTransactionLogger) Error() <-chan error {
	return f.errors
}

func (f *FileTransactionLogger) Run() {
	events := make(chan Event, 16)
	f.events = events

	errors := make(chan error, 1)
	f.errors = errors
	go func() {
		for e := range events {
			fmt.Printf("f")
			f.lastSequence++
			_, err := fmt.Fprintf(f.file,
				"%d\t%d\t%s\t%s\n",
				f.lastSequence, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
				return
			}
		}
	}()
}
func (f *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(f.file)
	outEvents := make(chan Event) // Небуферизованный канал событий
	outErrors := make(chan error, 1)
	go func() {
		var e Event
		defer close(outEvents)
		defer close(outErrors)
		for scanner.Scan() {

			line := scanner.Text()
			fmt.Printf("textLine: %s\n", line)
			if _, err := fmt.Sscanf(line, "%d\t%d\t%s\t%s",
				&e.Sequence, &e.EventType, &e.Key, &e.Value); err != nil {
				outErrors <- fmt.Errorf("input parse error: %w", err)
				return
			}
			if f.lastSequence >= e.Sequence {
				outErrors <- fmt.Errorf("transaction numbers out of sequence")
				return
			}
			f.lastSequence = e.Sequence
			fmt.Printf("value: %s\n", e.Value)
			outEvents <- e
		}
		if err := scanner.Err(); err != nil {
			outErrors <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}
	}()
	return outEvents, outErrors
}

func (f *FileTransactionLogger) Initialize(storage *Storage) error {
	var err error
	events, errorsLog := f.ReadEvents()
	e, ok := Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errorsLog:
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete:
				err = storage.Delete(e.Key)
			case EventPut:
				err = storage.Put(e.Key, string(e.Value))
			}
		}
	}
	f.Run()
	return err
}
func (f *FileTransactionLogger) Close() error {
	err := f.file.Close()
	if err != nil {
		return fmt.Errorf("error closing file %w", err)
	}
	return err
}
