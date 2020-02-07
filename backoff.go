/*

Package backoff implements a simple exponential backoff algorithms for retrying operations.

Use Do function for retrying operations that may fail.

Simple example:

	boff := backoff.New(
		backoff.Exponential(),
		backoff.MaxRetries(30),
		backoff.Notify(func(err error, next time.Duration, count uint) {
			log.Println(count, next,  err)
		}),
	)
	err := boff.Do(ctx, func(ctx) error {
		// if fail
		return err
		// if fail but not retry
		return backoff.Continue(err)
		// if ok
		return nil
	})
*/
package backoff

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"ciceroverneck.dev/backoff/internal"
)

const (
	defaultMaxRetries          = 0
	defaultInterval            = 500 * time.Millisecond
	defaultMaxInterval         = 60 * time.Second
	defaultMaxElapsedTime      = 0
	defaultRandomizationFactor = 0.5
	defaultMultiplier          = 1.5
)

var (
	ErrMaxRetries     = errors.New("")
	ErrMaxElapsedTime = errors.New("")
)


type Handle func(ctx context.Context) error
type NotifyHandle func(err error, duration time.Duration, count uint)

func Continue(err error) error {
	return &wrapContinueError{err: err}
}

type wrapContinueError struct {
	err error
}

func (this *wrapContinueError) Error() string {
	return this.err.Error()
}
func (this *wrapContinueError) Unwrap() error {
	return this.err
}

type Backoff struct {
	opts *internal.BackoffSettings
}

func New(opt ...Option) *Backoff {
	args := &internal.BackoffSettings{
		MaxRetries:          defaultMaxRetries,
		MaxElapsedTime:      defaultMaxElapsedTime,
		Multiplier:          defaultMultiplier,
		Interval:            defaultInterval,
		MaxInterval:         defaultMaxInterval,
		RandomizationFactor: defaultRandomizationFactor,
	}
	for _, setter := range opt {
		setter(args)
	}
	return &Backoff{
		opts: args,
	}
}

func (this *Backoff) Do(ctx context.Context, fn Handle) error {
	var start = time.Now()
	var currentInterval = this.opts.Interval
	var count uint
	var err error
	var errs []error

	for {
		if err = fn(ctx); err == nil {
			return nil
		}
		var errContinue *wrapContinueError
		if errors.As(err, &errContinue) {
			return errContinue.err
		}

		count++
		if this.opts.MaxRetries != 0 && count >= this.opts.MaxRetries {
			return ErrMaxRetries
		}
		if this.opts.MaxElapsedTime != 0 && time.Now().Sub(start) > this.opts.MaxElapsedTime {
			return ErrMaxElapsedTime
		}

		if float64(currentInterval) >= float64(this.opts.MaxInterval)/this.opts.Multiplier {
			currentInterval = this.opts.MaxInterval
		} else {
			if this.opts.Exponential {
				currentInterval = time.Duration(float64(currentInterval) * this.opts.Multiplier)
			}
		}
		currentInterval = getRandomValueFromInterval(this.opts.RandomizationFactor, currentInterval)

		if this.opts.Notify != nil {
			this.opts.Notify(err, currentInterval, count)
		}
		errs = append(errs, err)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(currentInterval):
		}
	}
}

func getRandomValueFromInterval(randomizationFactor float64, currentInterval time.Duration) time.Duration {
	delta := randomizationFactor * float64(currentInterval)
	minInterval := float64(currentInterval) - delta
	maxInterval := float64(currentInterval) + delta
	return time.Duration(minInterval + (rand.New(rand.NewSource(time.Now().UnixNano())).Float64() * (maxInterval - minInterval + 1)))
}
