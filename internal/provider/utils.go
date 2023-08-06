package provider

import (
	"time"
)

const maxRetries = 3
const retryDelay = 5 * time.Second

func retryOperation(operation func() error) error {
	var err error

	for retries := 0; retries < maxRetries; retries++ {
		err = operation()
		if err == nil {
			return nil
		}
		time.Sleep(retryDelay)
	}

	return err
}
