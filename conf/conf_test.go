package conf

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"golang.org/x/net/context"
)

// 这个测试会重置etcd里的设置，不建议跑
func TestWatchConf(t *testing.T) {
	t.Log(GetConfig())
	testConfig := APPConfig{
		StorageCfgs: make([]StorageConfig, 0),
		NodeCfg: NodeConfig{
			ListenAddr:   "test",
			Zone:         "test",
			Interval:     1,
			WriteTracing: 1,
			QueryTracing: 1,
		},
	}
	testConfig.StorageCfgs = append(testConfig.StorageCfgs, StorageConfig{
		URL: "http://localhost:8086",
	})
	data, err := json.Marshal(testConfig)
	if err != nil {
		t.Error("encode error", err)
	}
	resp, err := cli.Put(context.Background(), configPath, string(data))
	if err != nil {
		t.Error("put error", err)
	}
	t.Log(resp)
	time.Sleep(time.Second * 3)
	if !reflect.DeepEqual(GetConfig(), testConfig) {
		t.Error("watch error", GetConfig(), testConfig)
	}
}
