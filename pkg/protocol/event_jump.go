package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
)

type BigAirJumpEvent struct {
	JumpType      byte
	JumpNumber    uint16
	UnknownBytes1 [25]byte
	JumpTime      Time
	UnknownBytes2 [14]byte
}

func (e *BigAirJumpEvent) String() string {
	return fmt.Sprintf(
		`BigAirJumpEvent
  JumpType: %d
  JumpNumber: %d
  UnknownBytes1: % x
  JumpTime: %s
  UnknownBytes2: % x
`,
		e.JumpType, e.JumpNumber, e.UnknownBytes1, e.JumpTime, e.UnknownBytes2)
}

func readBigAirJumpEvent(reader io.Reader) (*BigAirJumpEvent, error) {
	result := &BigAirJumpEvent{}

	err := binary.Read(reader, binary.LittleEndian, &result.JumpType)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.LittleEndian, &result.JumpNumber)
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

	err = binary.Read(reader, binary.LittleEndian, &result.JumpTime)
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
