package service

type ClusterService struct {
}

func (s *ClusterService) Query(db string, q []byte) (results []byte, err error) {
	return
}

func (s *ClusterService) Write(db string, lp []byte) (err error) {
	return
}
