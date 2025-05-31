package protocol

import (
	"encoding/binary"
	"errors"
	"io"
)

type EventType byte

const (
	RecordingEventType EventType = 0x58
	BigAirJumpType     EventType = 0x44
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
	case BigAirJumpType:
		event, err := readBigAirJumpEvent(reader)
		if err != nil {
			return 0, nil, err
		}
		return BigAirJumpType, event, nil
	}

	return 0, nil, errors.New("unknown event type")
}
