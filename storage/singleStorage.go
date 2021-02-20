package storage

import (
	"fmt"
	"influxcluster/logging"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var (
	logger = logging.GetLogger("storage")
	DefaultCheckInterval = 1000
)

type SingleStorage struct {
	client    *http.Client
	transport http.Transport
	CheckInterval int
	BaseURL string
	Active    bool
	running   bool
}

func NewSingleStorage(baseURL string)(s *SingleStorage){
	s = &SingleStorage{
		client:    &http.Client{},
		BaseURL:   baseURL,
		CheckInterval: DefaultCheckInterval,
	}
	go s.checkActive()
	return
}

func (s *SingleStorage) Query(q string) (results []byte, err error){
	var req *http.Request
	req, err = http.NewRequest("GET",s.BaseURL + "/query?" + q,nil )
	if err != nil{
		return nil, err
	}
	resp, err := s.transport.RoundTrip(req)
	if err != nil{
		return nil, err
	}
	if resp.StatusCode != http.StatusOK{
		err = fmt.Errorf("response code error, resp:%v", resp)
		return
	}
	results, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("read body error: %s,the query is %s\n", err, q)
		return
	}
	logger.Infof("storage %s query %s",s.BaseURL,q)
	return
}

// lp, Line Protocol, contains the time series data that you want to store. Its components are measurement, tags, fields and timestamp.
func (s *SingleStorage) Write(db string, lp []byte)(err error){
	buf, err := gzipCompress(lp)
	if err != nil {
		logger.Error("compress error: ", err)
		return
	}

	logger.Infof("storage %s %s write %s",s.BaseURL,db,string(lp))
	err = s.WriteStream(db,buf, true)
	return
}

func (s *SingleStorage) WriteStream(db string, stream io.Reader, compressed bool) (err error) {
	q := url.Values{}
	q.Set("db", db)

	req, err := http.NewRequest("POST", s.BaseURL+"/write?"+q.Encode(), stream)
	if compressed {
		req.Header.Add("Content-Encoding", "gzip")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		logger.Error("http error: ", err)
		s.Active = false
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return
	}
	logger.Error("write status code: ", resp.StatusCode)

	respbuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("read body error: ", err)
		return
	}
	logger.Errorf("error response: %s\n", respbuf)

	// https://docs.influxdata.com/influxdb/v1.1/tools/api/#write
	switch resp.StatusCode {
	case 400:
		err = ErrBadRequest
	case 404:
		err = ErrNotFound
	default: // mostly tcp connection timeout
		logger.Errorf("status: %d", resp.StatusCode)
		err = ErrUnknown
	}
	return
}

func (s *SingleStorage) checkActive(){
	var err error
	for s.running {
		_, err = s.Ping()
		s.Active = err == nil
		time.Sleep(time.Millisecond * time.Duration(s.CheckInterval))
	}
	return
}

func (s *SingleStorage) Ping() (version string, err error){
	endpoint := s.BaseURL + "/ping"
	resp, err := s.client.Get(endpoint)
	if err != nil {
		logger.Info("network error: ", err)
		return
	}
	if resp != nil && resp.Body != nil{
		defer resp.Body.Close()
	}

	version = resp.Header.Get("X-Influxdb-Version")
	logger.Infof("ping %s, version:%s\n",endpoint,version)
	// Receive HTTP 204 No Content after writing the data, indicating that the writing was successful
	if resp.StatusCode == http.StatusNoContent {
		return
	}
	logger.Errorf("ping status code: %d, the base url is %s\n", resp.StatusCode, s.BaseURL)
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("read resp body error: ", err)
		return
	}
	logger.Errorf("error response: %s\n", respData)
	return
}

func (s *SingleStorage) Close() (err error){
	s.running = false
	s.transport.CloseIdleConnections()
	return
}