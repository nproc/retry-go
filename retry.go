package retry

import (
	"errors"
	"time"

	"github.com/nproc/errorgroup-go"
)

var (
	ErrTryFnNil = errors.New("TryFen can not be nil")
)

// TryFn is the function to try to execute.
//
// It receives as arguments:
// - the number of this attempt `[0..len(BackoffArray]`
// - the limit of executions `len(BackoffArray)`
type TryFn func(attempt, limit int) error

// BackoffArray an vector of interval to wait between each retry
type BackoffArray []time.Duration

// WithBackoffArray runs the TryFn with intervals from the given BackoffArray
//
// It is important to notice that the `TryFn` will run `len(BackoffArray) + 1`
// times
func WithBackoffArray(backoff BackoffArray, fn TryFn) error {
	return loop(cloneBackoffArray(backoff), fn)
}

// WithFixedInterval runs the TryFn with a BackoffArray created with the
// given `interval` repeated `repeat` times
//
// It is important to notice that the `TryFn` will run `repeat + 1` times
func WithFixedInterval(interval time.Duration, repeat int, fn TryFn) error {
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

func loop(backoff BackoffArray, fn TryFn) error {
	limit := len(backoff)

	if fn == nil {
		return ErrTryFnNil
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

func run(fn TryFn, attempt, limit int) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	err = fn(attempt, limit)

	return
}
