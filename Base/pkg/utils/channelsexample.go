package utils

import "fmt"

type Example struct {
	requests chan struct{}
	stopChan chan struct{}
}

func (e *Example) Run() {
	go func() {
		for {
			select {
			case request := <-e.requests:
				// Process the request
				fmt.Println(request)

			case <-e.stopChan:
				return
			}
		}
	}()
}
