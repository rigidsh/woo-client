package protocol

import "bytes"

type BufferPackageWriter struct {
	buff         *bytes.Buffer
	onPackageEnd func(checksum bool, data []byte)
}

func NewBufferPackageWriter(onPackageEnd func(checksum bool, data []byte)) *BufferPackageWriter {
	return &BufferPackageWriter{onPackageEnd: onPackageEnd}
}

func (w *BufferPackageWriter) PackageStart() {
	w.buff = bytes.NewBuffer(make([]byte, 0, 1024))
}

func (w *BufferPackageWriter) PackageEnd(checksum bool) {
	w.onPackageEnd(checksum, w.buff.Bytes())
}

func (w *BufferPackageWriter) Write(p []byte) (n int, err error) {
	return w.buff.Write(p)
}
