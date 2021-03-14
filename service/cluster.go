package service

import (
	"influxcluster/conf"
	"influxcluster/logging"
	"influxcluster/storage"
	"sync"
	"sync/atomic"
	"time"
)

var (
	logger = logging.GetLogger("cluster")
)

type ClusterService struct {
	storages []*storage.SingleStorage
	mu       sync.RWMutex
	stats    *Statistics
}

func (s *ClusterService) Query(db, q string) (results []byte, err error) {
	atomic.AddInt64(&s.stats.QueryRequests, 1)
	defer func(start time.Time) {
		atomic.AddInt64(&s.stats.QueryRequestDuration, time.Since(start).Nanoseconds())
	}(time.Now())

	for _, backStorage := range s.storages {
		if backStorage.Active {
			results, err = backStorage.Query(db, q)
			if err == nil {
				return
			}
		}
	}
	atomic.AddInt64(&s.stats.QueryRequestsFail, 1)
	return
}

func NewInfluxCluster(appcfg conf.APPConfig) (s *ClusterService) {
	s = &ClusterService{
		storages: make([]*storage.SingleStorage, 0),
		mu:       sync.RWMutex{},
		stats:    &Statistics{},
	}
	s.ReloadCfg(appcfg)
	return
}

func (s *ClusterService) ReloadCfg(appcfg conf.APPConfig) (err error) {
	var newStorages []*storage.SingleStorage
	for _, storageCfg := range appcfg.StorageCfgs {
		newStorages = append(newStorages, storage.NewSingleStorage(storageCfg))
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, backStorage := range s.storages {
		err = backStorage.Close()
		if err != nil {
			logger.Infof("fail in close backend %s", backStorage.BaseURL)
		}
	}
	s.storages = newStorages
	return
}

func (s *ClusterService) Write(db string, lp []byte) (err error) {
	atomic.AddInt64(&s.stats.WriteRequests, 1)
	defer func(start time.Time) {
		atomic.AddInt64(&s.stats.WriteRequestDuration, time.Since(start).Nanoseconds())
	}(time.Now())

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, backStorage := range s.storages {
		err = backStorage.Write(db, lp)
		if err != nil {
			logger.Infof("write error: %s\n", err)
			atomic.AddInt64(&s.stats.WriteRequestsFail, 1)
		}
	}

	return
}

func (s *ClusterService) Close() (err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, backStorage := range s.storages {
		err = backStorage.Close()
		if err != nil {
			logger.Infof("fail in close backend %s", backStorage.BaseURL)
		}
	}
	return
}
