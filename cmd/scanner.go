package cmd

type ScanResult struct {
	ConcurrencyLimit int
	Url              string
	PClientResponses []PClientResponse
	RequestCount     int
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

	for i := 0; i < requestCount; i++ {
		go func(i int) {
			semaphoreChan <- struct{}{}
			pclientResponse, _ := scanner.PClient.call(url)
			responseChan <- pclientResponse
			<-semaphoreChan
		}(i)
	}

	var responses []PClientResponse

	for {
		response := <-responseChan
		responses = append(responses, *response)
		if len(responses) == requestCount {
			break
		}
	}
	return &ScanResult{
		PClientResponses: responses,
		ConcurrencyLimit: concurrencyLimit,
		RequestCount:     requestCount,
		Url:              url,
	}, nil
}
