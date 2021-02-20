package storage

import (
	"net/url"
	"testing"
)

func TestPing(t *testing.T) {
	s := NewSingleStorage("http://localhost:8086")
	defer s.Close()

	version, err := s.Ping()
	if err != nil {
		t.Errorf("error: %s", err)
		return
	}
	if version == "" {
		t.Errorf("empty version")
	}
	return
}

func TestQuery(t *testing.T){
	s := NewSingleStorage("http://localhost:8086")
	defer s.Close()
	q := make(url.Values, 1)
	q.Set("db", "mydb")
	q.Set("q", "SELECT * FROM \"cpu_load_short\"")
	results, err := s.Query(q.Encode())
	if err != nil {
		t.Errorf("error: %s", err)
		return
	}
	t.Log(string(results))
}

func TestWrite(t *testing.T) {
	s := NewSingleStorage("http://localhost:8086")
	defer s.Close()

	err := s.Write("mydb",[]byte("cpu,host=server01,region=uswest value=1 1434055562000000000\ncpu value=3,value2=4 1434055562000010000"))
	if err != nil {
		t.Errorf("error: %s", err)
		return
	}
}