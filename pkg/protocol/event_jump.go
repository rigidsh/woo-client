package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
)

type JumpEvent struct {
	JumpType      byte
	JumpNumber    uint16
	UnknownBytes1 [10]byte
	JumpHeight    uint16
	UnknownBytes2 [13]byte
	JumpTime      Time
	UnknownBytes3 [14]byte
}

func (e *JumpEvent) String() string {
	return fmt.Sprintf(
		`JumpEvent
  JumpType: %d
  JumpNumber: %d
  UnknownBytes1: % x
  JumpHeight: %d
  UnknownBytes2: % x
  JumpTime: %s
  UnknownBytes3: % x
`,
		e.JumpType, e.JumpNumber, e.UnknownBytes1, e.JumpHeight, e.UnknownBytes2, e.JumpTime, e.UnknownBytes3)
}

func readJumpEvent(reader io.Reader) (*JumpEvent, error) {
	result := &JumpEvent{}

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

	err = binary.Read(reader, binary.LittleEndian, &result.JumpHeight)
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

	err = binary.Read(reader, binary.LittleEndian, &result.JumpTime)
	if err != nil {
		return nil, err
	}

	n, err = reader.Read(result.UnknownBytes3[:])
	if err != nil {
		return nil, err
	}
	if n != len(result.UnknownBytes3) {
		return nil, io.ErrUnexpectedEOF
	}

	return result, nil
}
