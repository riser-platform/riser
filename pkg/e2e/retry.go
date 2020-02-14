package e2e

import (
	"time"

	"github.com/pkg/errors"
)

const defaultRetryDelay = 1 * time.Second
const defaultMaxRetries = 30

type RetryFunc func() (abort bool, err error)

func Retry(fn RetryFunc) error {
	var err error
	for i := 1; i <= defaultMaxRetries; i++ {
		var abort bool
		abort, err = fn()
		if abort {
			return err
		}
		<-time.After(defaultRetryDelay)
	}

	return errors.Wrap(err, "Retries exhausted")
}
