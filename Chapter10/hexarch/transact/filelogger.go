package transact

import (
	"bufio"
	"fmt"
	"github.com/renanvicente/CloudNativeGo/Chapter10/hexarch/core"
	"os"
)

type FileTransactionLogger struct {
	events       chan<- core.Event // Write-only channel for sending events
	errors       <-chan error      // Read-only channel for receiving errors
	lastSequence uint64            // The last used event sequence number
	file         *os.File          // The location of the transaction log

}

func (l *FileTransactionLogger) WritePut(key, value string) {
	l.events <- core.Event{
		EventType: core.EventPut,
		Key:       key,
		Value:     value,
	}

}

func (l *FileTransactionLogger) WriteDelete(key string) {
	l.events <- core.Event{
		EventType: core.EventDelete,
		Key:       key,
	}
}

func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *FileTransactionLogger) Run() {
	events := make(chan core.Event, 16) // Make an events channel
	l.events = events

	errors := make(chan error, 1) // Make an errors channel
	l.errors = errors

	go func() {
		for e := range events { // Retrieve the next Event

			l.lastSequence++ // Increment sequence number

			_, err := fmt.Fprintf( // Write the event to the log
				l.file,
				"%d\t%d\t%s\t%s\n",
				l.lastSequence, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
				return
			}

		}
	}()

}

func (l *FileTransactionLogger) ReadEvents() (<-chan core.Event, <-chan error) {
	scanner := bufio.NewScanner(l.file) // Create a Scanner for l.file
	outEvent := make(chan core.Event)   // An unbuffered Event channel
	outError := make(chan error, 1)     // A buffered error channel

	go func() {
		var e core.Event
		defer close(outEvent) // Close the channels when the
		defer close(outError) // goroutine ends

		for scanner.Scan() {
			line := scanner.Text()
			fmt.Sscanf(line, "%d\t%d\t%s\t%s", &e.Sequence, &e.EventType, &e.Key, &e.Value)

			//if _, err := fmt.Sscanf(line, "%d\t%d\t%s\t%s", &e.Sequence, &e.EventType, &e.Key, &e.Value); err != nil {
			//	outError <- fmt.Errorf("input parse error: %w", err)
			//	return
			//}
			//if err != nil {
			//	outError <- fmt.Errorf("input parse error: %w", err)
			//	return
			//}
			// Sanity check! Are the sequence numbers in increasing order?
			if l.lastSequence >= e.Sequence {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}
			l.lastSequence = e.Sequence
			outEvent <- e

		}
		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}

	}()
	return outEvent, outError
}
func NewFileTransactionLogger(filename string) (core.TransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
	}
	return &FileTransactionLogger{file: file}, nil

}
