package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
)

type RecordingEvent struct {
	Start          bool
	UnknownBytes1  [1]byte
	NumberOfEvents uint16
	Counter        uint32
	UnknownBytes2  [4]byte
}

func (e *RecordingEvent) String() string {
	return fmt.Sprintf(
		`RecordingEvent
  Start: %t
  UnknownBytes1: % x
  NumberOfEvents: %d
  Counter: %d
  UnknownBytes2: % x
`,
		e.Start, e.UnknownBytes1, e.NumberOfEvents, e.Counter, e.UnknownBytes2)
}

func readRecordingEvent(reader io.Reader) (*RecordingEvent, error) {
	result := &RecordingEvent{}

	err := binary.Read(reader, binary.LittleEndian, &result.Start)
	if err != nil {
		return nil, err
	}

	n, err := reader.Read(result.UnknownBytes1[:])
	if err != nil {
		return nil, err
	}
	if n != len(result.UnknownBytes1) {
		return nil, io.ErrUnexpectedEOF
	}

	err = binary.Read(reader, binary.LittleEndian, &result.NumberOfEvents)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.LittleEndian, &result.Counter)
	if err != nil {
		return nil, err
	}

	n, err = reader.Read(result.UnknownBytes2[:])
	if err != nil {
		return nil, err
	}
	if n != len(result.UnknownBytes2) {
		return nil, io.ErrUnexpectedEOF
	}

	return result, nil
}
