package storage

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io"
	"io/ioutil"
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

func gzipUnCompress(p []byte) (data []byte, err error) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, p)
	zip, err := gzip.NewReader(buf)
	if err != nil {
		return
	}
	data, err = ioutil.ReadAll(zip)
	if err != nil {
		return
	}
	err = zip.Close()
	return
}
