package cmd

import (
	"encoding/json"
	"time"
)

type ScanResult struct {
	ConcurrencyLimit          int
	Url                       string
	PClientResponses          []PClientResponse
	RequestCount              int
	DurationNanoSeconds       time.Duration
	DurationMilliSeconds      int64
	RequestResponseEvaluation RequestResponseEvaluation
}

type IntermediateScanResult struct {
	DurationNanoSeconds       time.Duration
	DurationMilliSeconds      int64
	RequestResponseEvaluation RequestResponseEvaluation
}

type RequestResponseEvaluation struct {
	SuccessCount       int
	ClientFailureCount int
	ServerFailureCount int
}

type Scanner struct {
	PClient *PClient
	Logging *Logging
}

func (scanner Scanner) scan(url string, requestCount int, concurrencyLimit int) (*ScanResult, error) {

	semaphoreChan := make(chan struct{}, concurrencyLimit)
	responseChan := make(chan PClientResponse)

	defer func() {
		close(semaphoreChan)
		close(responseChan)
	}()

	scanner.Logging.DebugLog.Printf("start scanning")
	startTime := time.Now()
	for i := 0; i < requestCount; i++ {
		go func(i int) {
			semaphoreChan <- struct{}{}
			scanner.Logging.DebugLog.Printf("Call PClient")
			pclientResponse := scanner.PClient.call(url)
			responseChan <- pclientResponse
			<-semaphoreChan
		}(i)
	}

	var responses []PClientResponse
	for {
		response := <-responseChan
		responses = append(responses, response)

		logResponse(response, scanner)
		logIntermediateScanResult(startTime, responses, scanner)

		if len(responses) == requestCount {
			break
		}
	}

	scanner.Logging.DebugLog.Printf("finished scanning")

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	scanResult := ScanResult{
		ConcurrencyLimit:          concurrencyLimit,
		Url:                       url,
		PClientResponses:          responses,
		RequestCount:              requestCount,
		DurationNanoSeconds:       duration,
		DurationMilliSeconds:      duration.Milliseconds(),
		RequestResponseEvaluation: evaluateResponses(responses),
	}
	return &scanResult, nil
}

func logIntermediateScanResult(startTime time.Time, responses []PClientResponse, scanner Scanner) {
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	intermediateScanResult := IntermediateScanResult{
		RequestResponseEvaluation: evaluateResponses(responses),
		DurationNanoSeconds:       duration,
		DurationMilliSeconds:      duration.Milliseconds(),
	}
	intermediateScanResultJson, err := json.Marshal(intermediateScanResult)
	if err != nil {
		scanner.Logging.ErrorLog.Printf("%v", err)
	} else {
		scanner.Logging.Result.Printf(string(intermediateScanResultJson))
	}
}

func logResponse(response PClientResponse, scanner Scanner) {
	responseJson, err := json.Marshal(response)
	if err != nil {
		scanner.Logging.ErrorLog.Printf("%v", err)
	} else {
		scanner.Logging.Result.Printf(string(responseJson))
	}
}

func evaluateResponses(responses []PClientResponse) RequestResponseEvaluation {
	successCount := 0
	clientFailureCount := 0
	serverFailureCount := 0
	for _, response := range responses {
		is2xx := checkStatusCodeIs2xx(response.StatusCode)
		if response.Error == nil && is2xx {
			successCount++
		} else if response.Error != nil {
			clientFailureCount++
		} else if !is2xx {
			serverFailureCount++
		}
	}
	return RequestResponseEvaluation{successCount, clientFailureCount, serverFailureCount}
}

func checkStatusCodeIs2xx(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}
