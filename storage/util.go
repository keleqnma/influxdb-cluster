package storage

import (
	"bytes"
	"compress/gzip"
	"io"
)

func gzipCompress(p []byte) (buf *bytes.Buffer, err error) {
	buf = &bytes.Buffer{}
	zip := gzip.NewWriter(buf)
	n, err := zip.Write(p)
	if err != nil {
		return
	}
	if n != len(p) {
		err = io.ErrShortWrite
		return
	}
	err = zip.Close()
	return
}
