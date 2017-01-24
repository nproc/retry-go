package retry

import (
	"errors"
	"time"

	"github.com/txgruppi/errorgroup-go"
)

var (
	// ErrTryFuncNil is returned when the TryFunc is nil
	ErrTryFuncNil = errors.New("TryFunc can not be nil")
)

// TryFunc is the function to try to execute.
//
// It receives as arguments:
// - the number of this attempt '[0..len(BackoffArray]'
// - the limit of executions 'len(BackoffArray) + 1'
type TryFunc func(attempt, limit int) error

// BackoffArray an vector of interval to wait between each retry
type BackoffArray []time.Duration

// WithBackoffArray runs the TryFunc with intervals from the given BackoffArray
//
// It is important to notice that the TryFunc will run 'len(BackoffArray) + 1'
// times
func WithBackoffArray(backoff BackoffArray, fn TryFunc) error {
	return loop(cloneBackoffArray(backoff), fn)
}

// WithFixedInterval runs the TryFunc with a BackoffArray created with the
// given 'interval' repeated 'repeat' times
//
// It is important to notice that the TryFunc will run 'repeat + 1' times
func WithFixedInterval(interval time.Duration, repeat int, fn TryFunc) error {
	backoff := make([]time.Duration, repeat)
	for i := 0; i < repeat; i++ {
		backoff[i] = interval
	}
	return loop(backoff, fn)
}

func cloneBackoffArray(backoff BackoffArray) BackoffArray {
	clone := make([]time.Duration, len(backoff))
	for i, v := range backoff {
		clone[i] = v
	}
	return clone
}

func loop(backoff BackoffArray, fn TryFunc) error {
	limit := len(backoff) + 1

	if fn == nil {
		return ErrTryFuncNil
	}

	errs := []error{}
	attempt := 0
	for {
		err := run(fn, attempt, limit)
		if err == nil {
			return nil
		}
		errs = append(errs, err)
		if len(backoff) == 0 {
			break
		}
		time.Sleep(backoff[0])
		backoff = backoff[1:]
		attempt++
	}

	return errorgroup.New(errs)
}

func run(fn TryFunc, attempt, limit int) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	err = fn(attempt, limit)

	return
}
