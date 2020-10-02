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
	StartTime time.Time
	EndTime	  time.Time
	StatusCode int
}


func (pclient PClient) call(url string) (*PClientResponse, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: time.Second * 1}
	startTime := time.Now()
	resp, err := client.Do(req)
	endTime := time.Now()

	if err != nil {
		return nil, err
	}

	pclientResponse := &PClientResponse{
		StatusCode: resp.StatusCode,
		StartTime: startTime,
		EndTime: endTime,
	}

	return pclientResponse, nil
}