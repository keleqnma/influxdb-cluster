package storage

import (
	"influxcluster/conf"
	"testing"
)

func GetTestCfg() conf.StorageConfig {
	return conf.StorageConfig{
		URL:           "http://localhost:8086",
		CheckInterval: 1000,
	}
}

func TestPing(t *testing.T) {
	s := NewSingleStorage(GetTestCfg())
	defer s.Close()

	version, err := s.Ping()
	if err != nil {
		t.Errorf("error: %s", err)
		return
	}
	if version == "" {
		t.Errorf("empty version")
	}
}

func TestQuery(t *testing.T) {
	s := NewSingleStorage(GetTestCfg())
	defer s.Close()
	results, err := s.Query("mydb", "SELECT * FROM \"cpu_load_short\"")
	if err != nil {
		t.Errorf("error: %s", err)
		return
	}
	t.Log(string(results))
}

func TestWrite(t *testing.T) {
	s := NewSingleStorage(GetTestCfg())
	defer s.Close()

	err := s.Write("mydb", []byte("cpu,host=server01,region=uswest value=1 1434055562000000000\ncpu value=3,value2=4 1434055562000010000"))
	if err != nil {
		t.Errorf("error: %s", err)
		return
	}
}

func TestQueryAfterWrite(t *testing.T) {
	s := NewSingleStorage(GetTestCfg())
	defer s.Close()

	// write
	err := s.Write("mydb", []byte("cpu,host=server01,region=uswest value=1 1434055562000000000\ncpu value=3,value2=4 1434055562000010000"))
	if err != nil {
		t.Errorf("error: %s", err)
		return
	}

	// query
	results, err := s.Query("mydb", "SELECT * FROM \"cpu_load_short\"")
	if err != nil {
		t.Errorf("error: %s", err)
		return
	}
	t.Log(string(results))
}
