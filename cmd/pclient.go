package cmd

import (
	"net/http"
	"time"
)

type PClient struct {
	Logging *Logging
}

func NewPClient(logging Logging) *PClient {
	return &PClient{
		Logging: &logging,
	}
}

type PClientResponse struct {
	StartTime  time.Time
	EndTime    time.Time
	StatusCode int
	Duration   time.Duration
	Error	   error
}

func NewPClientResponse(startTime time.Time, endTime time.Time, statusCode int, err error) PClientResponse {

	duration :=  endTime.Sub(startTime)

	return PClientResponse{
		StatusCode: statusCode,
		StartTime:  startTime,
		EndTime:    endTime,
		Duration: 	duration,
		Error:		err,
	}
}

func (pclient PClient) call(url string) PClientResponse {
	req, err := http.NewRequest("GET", url, nil)

	pclient.Logging.DebugLog.Printf("call %s", url)

	if err != nil {
		return PClientResponse{Error: err}
	}

	client := &http.Client{Timeout: time.Second * 1}
	startTime := time.Now()
	resp, err := client.Do(req)
	endTime := time.Now()

	var pclientResponse PClientResponse

	if resp != nil {
		pclientResponse = NewPClientResponse(startTime, endTime, resp.StatusCode, err)
	}  else {
		pclientResponse = NewPClientResponse(startTime, endTime, 0, err)
	}

	return pclientResponse
}
