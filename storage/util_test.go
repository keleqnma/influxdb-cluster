package storage

import (
	"bytes"
	"testing"
)

func TestZipUnzip(t *testing.T) {
	s := []byte("testingtesting")
	buf, err := gzipCompress(s)
	if err != nil {
		t.Errorf("zip error: %v", err)
	}
	zipS := buf.Bytes()
	unzipS, err := gzipUnCompress(zipS)
	if err != nil {
		t.Errorf("unZip error: %v", err)
	}
	if !bytes.Equal(s, unzipS) {
		t.Error("error", string(s), string(zipS), string(unzipS))
	}
}
