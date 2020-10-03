package cmd

import "time"

type ScanResult struct {
	ConcurrencyLimit    int
	Url                 string
	PClientResponses    []PClientResponse
	RequestCount        int
	DurationNanoSeconds time.Duration
	DurationMilliSeconds int64
}

type Scanner struct {
	PClient *PClient
	Logging *Logging
}

func (scanner Scanner) scan(url string, requestCount int, concurrencyLimit int) (*ScanResult, error) {

	semaphoreChan := make(chan struct{}, concurrencyLimit)
	responseChan := make(chan *PClientResponse)

	defer func() {
		close(semaphoreChan)
		close(responseChan)
	}()


	startTime := time.Now()
	for i := 0; i < requestCount; i++ {
		go func(i int) {
			semaphoreChan <- struct{}{}
			pclientResponse, _ := scanner.PClient.call(url)
			responseChan <- pclientResponse
			<-semaphoreChan
		}(i)
	}
	endTime := time.Now()

	duration := endTime.Sub(startTime)

	var responses []PClientResponse
	for {
		response := <-responseChan
		responses = append(responses, *response)
		if len(responses) == requestCount {
			break
		}
	}

	scanResult := ScanResult{
		ConcurrencyLimit: concurrencyLimit,
		Url: url,
		PClientResponses: responses,
		RequestCount: requestCount,
		DurationNanoSeconds: duration,
		DurationMilliSeconds: duration.Milliseconds(),
	}
	return &scanResult , nil
}
