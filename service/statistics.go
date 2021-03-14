package service

type Statistics struct {
	QueryRequests        int64
	QueryRequestsFail    int64
	WriteRequests        int64
	WriteRequestsFail    int64
	PingRequests         int64
	PingRequestsFail     int64
	PointsWritten        int64
	PointsWrittenFail    int64
	WriteRequestDuration int64
	QueryRequestDuration int64
}

func (s *Statistics) Flush() {
	s.QueryRequests = 0
	s.QueryRequestsFail = 0
	s.WriteRequests = 0
	s.WriteRequestsFail = 0
	s.PingRequests = 0
	s.PingRequestsFail = 0
	s.PointsWritten = 0
	s.PointsWrittenFail = 0
	s.WriteRequestDuration = 0
	s.QueryRequestDuration = 0
}
