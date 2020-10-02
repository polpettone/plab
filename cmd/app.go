package cmd



type Application struct {
	Logging *Logging
	PClient *PClient
}


func NewApplication(logging *Logging, pclient *PClient) *Application {
	return &Application{
		Logging: logging,
		PClient: pclient,
	}
}

