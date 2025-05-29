package protocol

import (
	"errors"
	"io"
)

type state int

const (
	unknownState      state = 0
	inPacketState     state = 1
	escapeSymbolState state = 2
)

type Package []byte

type PackageWriter interface {
	io.Writer
	PackageStart()
	PackageEnd(checksum bool)
}

type PackageDecoder struct {
	target PackageWriter

	buff         oneByteBuffer
	checksum     byte
	currentState state
}

func NewPackageDecoder(target PackageWriter) *PackageDecoder {
	return &PackageDecoder{target: target}
}

func (decoder *PackageDecoder) Write(srcData []byte) (n int, err error) {
	data := newSliceWithBuffer(decoder.buff, srcData)
	decoder.buff = oneByteBuffer{}

	startPosition := 0

	for startPosition < data.Len() {
		symbol, symbolPosition := findNextControlSymbol(data.SliceFrom(startPosition), decoder.currentState == escapeSymbolState)

		if symbolPosition == -1 {
			if decoder.currentState == inPacketState || decoder.currentState == escapeSymbolState {
				writtenBytes, err := decoder.writeToTarget(data.Slice(startPosition, data.Len()-1))
				if err != nil {
					return 0, err
				}

				if writtenBytes != data.Len()-1-startPosition {
					return symbolPosition + writtenBytes, nil
				}
				if writtenBytes != 0 {
					decoder.currentState = inPacketState
				}

				decoder.buff.data = data.Get(data.Len() - 1)
				decoder.buff.hasData = true

				return data.ToDataIndex(data.Len()), nil
			}
		}

		if decoder.currentState == escapeSymbolState {
			decoder.currentState = inPacketState
		}

		switch symbol {
		case packageStart:
			if decoder.currentState != unknownState {
				return 0, errors.New("invalid control symbol")
			}
			startPosition += symbolPosition + 1
			decoder.target.PackageStart()
			decoder.currentState = inPacketState
		case packageEnd:
			if decoder.currentState != inPacketState {
				return 0, errors.New("invalid control symbol")
			}

			writtenBytes, err := decoder.writeToTarget(data.Slice(startPosition, startPosition+symbolPosition-1))
			if err != nil {
				return 0, err
			}

			if writtenBytes != symbolPosition-1 {
				return data.ToDataIndex(startPosition + writtenBytes), nil
			}

			decoder.target.PackageEnd(decoder.doChecksum(data.Get(startPosition + symbolPosition - 1)))
			decoder.currentState = unknownState
			decoder.checksum = 0

			startPosition += symbolPosition + 1
		case escapeSymbol:
			if decoder.currentState != inPacketState {
				return 0, errors.New("invalid control symbol")
			}

			writtenBytes, err := decoder.writeToTarget(data.Slice(startPosition, startPosition+symbolPosition))
			if err != nil {
				return 0, err
			}

			if writtenBytes != symbolPosition {
				return data.ToDataIndex(symbolPosition + writtenBytes), nil
			}
			decoder.currentState = escapeSymbolState

			startPosition += symbolPosition + 1
		}
	}

	return data.ToDataIndex(data.Len()), nil
}

func (decoder *PackageDecoder) writeToTarget(data sliceWithBuffer) (int, error) {
	writtenBytes, err := data.WriteTo(decoder.target)
	if err != nil {
		return 0, err
	}

	for i := 0; i < writtenBytes; i++ {
		decoder.checksum += data.Get(i)
	}

	return writtenBytes, nil
}

func (decoder *PackageDecoder) doChecksum(checksum byte) bool {
	return decoder.checksum+checksum == 0
}
