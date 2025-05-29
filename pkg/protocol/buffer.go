package protocol

import "io"

type oneByteBuffer struct {
	hasData bool
	data    byte
}

type sliceWithBuffer struct {
	buffer oneByteBuffer
	data   []byte
}

func newSliceWithBuffer(buffer oneByteBuffer, data []byte) sliceWithBuffer {
	return sliceWithBuffer{buffer: buffer, data: data}
}

func (s sliceWithBuffer) Len() int {
	if s.buffer.hasData {
		return len(s.data) + 1
	}
	return len(s.data)
}

func (s sliceWithBuffer) Get(i int) byte {
	if s.buffer.hasData {
		if i == 0 {
			return s.buffer.data
		}
		return s.data[i-1]
	}

	return s.data[i]
}

func (s sliceWithBuffer) FromDataIndex(index int) int {
	if s.buffer.hasData {
		return index + 1
	}
	return index
}

func (s sliceWithBuffer) ToDataIndex(index int) int {
	if s.buffer.hasData {
		return index - 1
	}
	return index
}

func (s sliceWithBuffer) SliceFrom(position int) sliceWithBuffer {
	return s.Slice(position, s.Len())
}

func (s sliceWithBuffer) Slice(start, end int) sliceWithBuffer {
	if end == start {
		return sliceWithBuffer{
			buffer: oneByteBuffer{},
			data:   []byte{},
		}
	}

	if start == 0 {
		if s.buffer.hasData {
			return sliceWithBuffer{s.buffer, s.data[0 : end-1]}
		} else {
			return sliceWithBuffer{s.buffer, s.data[0:end]}
		}
	}

	if s.buffer.hasData {
		start -= 1
		end -= 1
	}

	return sliceWithBuffer{data: s.data[start:end]}
}

func (s sliceWithBuffer) WriteTo(target io.Writer) (int, error) {
	if s.buffer.hasData {
		written, err := target.Write([]byte{s.buffer.data})
		if err != nil {
			return 0, err
		}
		if written != 1 {
			return 0, nil
		}
	}

	if s.Len() == 0 {
		return 0, nil
	}

	written, err := target.Write(s.data)
	if err != nil {
		return 0, err
	}

	if s.buffer.hasData {
		return written + 1, nil
	} else {
		return written, nil
	}
}
