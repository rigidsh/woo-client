package protocol

import (
	"bytes"
	"slices"
	"testing"
)

type FooPackageWriter struct {
	buff          *bytes.Buffer
	onPackageEnd  func(checksum bool, data []byte)
	limitOnTarget int
}

func NewFooPackageWriter(onPackageEnd func(checksum bool, data []byte), limitOnTarget int) *FooPackageWriter {
	return &FooPackageWriter{onPackageEnd: onPackageEnd, limitOnTarget: limitOnTarget}
}

func (w *FooPackageWriter) PackageStart() {
	w.buff = bytes.NewBuffer(make([]byte, 0, 1024))
}

func (w *FooPackageWriter) PackageEnd(checksum bool) {
	w.onPackageEnd(checksum, w.buff.Bytes())
}

func (w *FooPackageWriter) Write(p []byte) (n int, err error) {
	if w.limitOnTarget == -1 {
		return w.buff.Write(p)
	}

	return w.buff.Write(p[:w.limitOnTarget])
}

func TestNewPackageDecoder_empty_package(t *testing.T) {
	successTest(src([]byte{0xD1, 0x00, 0xDF}), []byte{}, -1, t)
}

func TestNewPackageDecoder_simple(t *testing.T) {
	successTest(src([]byte{0xD1, 0xFF, 0x01, 0xDF}), []byte{0xFF}, -1, t)
	successTest(src([]byte{0xD1, 0xFF, 0x01, 0xDF, 0xD1, 0xFF, 0x01, 0xDF}), []byte{0xFF}, -1, t)
	successTest(src([]byte{0xD1}, []byte{0x00}, []byte{0x00}, []byte{0xDF}), []byte{0x00}, -1, t)
	successTest(src([]byte{0xD1, 0x00}, []byte{0x00, 0xDF}), []byte{0x00}, -1, t)
}

func TestNewPackageDecoder_with_escape(t *testing.T) {
	successTest(src([]byte{0xD1, 0xDE, 0xD1, 0x2F, 0xDF}), []byte{0xD1}, -1, t)
	successTest(src([]byte{0xD1, 0xDE, 0xDE, 0x22, 0xDF}), []byte{0xDE}, -1, t)
	successTest(src([]byte{0xD1}, []byte{0xDE}, []byte{0xD1}, []byte{0x2F}, []byte{0xDF}), []byte{0xD1}, -1, t)
	successTest(src([]byte{0xD1, 0xDE}, []byte{0xD1, 0x2F}, []byte{0xDF}), []byte{0xD1}, -1, t)
}

func TestNewPackageDecoder_slow_target(t *testing.T) {
	successTest(src([]byte{0xD1, 0x00, 0x00, 0xDF}), []byte{0xD1}, 1, t)
	successTest(src([]byte{0xD1, 0x00, 0x00, 0xDE, 0xDE, 0x22, 0xDF}), []byte{0xDE}, -1, t)
	successTest(src([]byte{0xD1, 0x01, 0x02, 0x03, 0x04, 0xF6, 0xDF}), []byte{0xD1}, 1, t)
}

type srcParam [][]byte

func src(batch ...[]byte) srcParam {
	return batch
}

func successTest(srcData srcParam, expectedPackageData []byte, limitOnTarget int, t *testing.T) {
	hasPackage := false
	target := NewFooPackageWriter(func(checksum bool, data []byte) {
		hasPackage = true
		if !checksum {
			t.Errorf("checksum is not correct")
		}

		slices.Equal(data, expectedPackageData)
	}, limitOnTarget)

	decoder := NewPackageDecoder(target)
	for _, data := range srcData {
		toWrite := len(data)
		for toWrite > 0 {
			written, err := decoder.Write(data[len(data)-toWrite:])
			if err != nil {
				t.Errorf("error on write: %s", err.Error())
				return
			}
			if limitOnTarget == -1 && written != len(data) {
				t.Errorf("written bytes is not correct")
			}
			toWrite -= written
		}
	}

	if !hasPackage {
		t.Errorf("no pacakge")
	}
}
