package protocol

import (
	"encoding/binary"
	"errors"
	"io"
)

type EventType byte

const (
	RecordingEventType EventType = 0x58
	JumpEventType      EventType = 0x44
)

func ReadEvent(reader io.Reader) (EventType, interface{}, error) {
	var eventType EventType
	err := binary.Read(reader, binary.LittleEndian, &eventType)
	if err != nil {
		return 0, nil, err
	}
	switch eventType {
	case RecordingEventType:
		event, err := readRecordingEvent(reader)
		if err != nil {
			return 0, nil, err
		}
		return RecordingEventType, event, nil
	case JumpEventType:
		event, err := readJumpEvent(reader)
		if err != nil {
			return 0, nil, err
		}
		return JumpEventType, event, nil
	}

	return 0, nil, errors.New("unknown event type")
}
