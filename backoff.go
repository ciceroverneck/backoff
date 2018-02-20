/*

Package backoff implements a simple exponential backoff algorithms for retrying operations.

Use Do function for retrying operations that may fail.

Simple example:

	boff := backoff.New(
		backoff.Exponential(),
		backoff.MaxRetries(30),
		backoff.FailCallback(func(err error, count uint) {
			log.Println(count, err)
		}),
	)
	err := boff.Do(func() error {
		// if fail
		return err
		// if ok
		return nil
	})

*/
package backoff

import (
	"errors"
	"math/rand"
	"time"
)

const (
	defaultMaxRetries          = 100
	defaultInterval            = 500 * time.Millisecond
	defaultMaxInterval         = 60 * time.Second
	defaultMaxElapsedTime      = 15 * time.Minute
	defaultRandomizationFactor = 0.5
	defaultMultiplier          = 1.5
)

var (
	ErrMaxRetries     = errors.New("")
	ErrMaxElapsedTime = errors.New("")
)

type options struct {
	maxRetries          uint
	exponential         bool
	failCallback        FailHandle
	multiplier          float64
	randomizationFactor float64
	maxElapsedTime      time.Duration
	interval            time.Duration
	maxInterval         time.Duration
}

type Option func(*options)

type Handle func() error
type FailHandle func(err error, count uint)

func MaxRetries(v uint) Option {
	return func(args *options) {
		args.maxRetries = v
	}
}

func Exponential() Option {
	return func(args *options) {
		args.exponential = true
	}
}

func FailCallback(fn FailHandle) Option {
	return func(args *options) {
		args.failCallback = fn
	}
}

func MaxElapsedTime(duration time.Duration) Option {
	return func(args *options) {
		args.maxElapsedTime = duration
	}
}

func RandomizationFactor(randomizationFactor float64) Option {
	return func(args *options) {
		args.randomizationFactor = randomizationFactor
	}
}

func Multiplier(multiplier float64) Option {
	return func(args *options) {
		args.multiplier = multiplier
	}
}

func MaxInterval(duration time.Duration) Option {
	return func(args *options) {
		args.maxInterval = duration
	}
}

func Interval(duration time.Duration) Option {
	return func(args *options) {
		args.interval = duration
	}
}

type Backoff struct {
	opts *options
}

func New(opt ...Option) *Backoff {
	args := &options{
		maxRetries:          defaultMaxRetries,
		maxElapsedTime:      defaultMaxElapsedTime,
		multiplier:          defaultMultiplier,
		interval:            defaultInterval,
		maxInterval:         defaultMaxInterval,
		randomizationFactor: defaultRandomizationFactor,
	}
	for _, setter := range opt {
		setter(args)
	}
	return &Backoff{
		opts: args,
	}
}

func (b *Backoff) Do(fn Handle) error {
	start := time.Now()
	currentInterval := b.opts.interval
	var count uint
	for {
		if err := fn(); err != nil {
			if b.opts.failCallback != nil {
				b.opts.failCallback(err, count)
			}
			if count > b.opts.maxRetries {
				return ErrMaxRetries
			}
			if b.opts.maxElapsedTime != 0 && time.Now().Sub(start) > b.opts.maxElapsedTime {
				return ErrMaxElapsedTime
			}
			count++
			if float64(currentInterval) >= float64(b.opts.maxInterval)/b.opts.multiplier {
				currentInterval = b.opts.maxInterval
			} else {
				if b.opts.exponential {
					currentInterval = time.Duration(float64(currentInterval) * b.opts.multiplier)
				}
			}

			currentInterval = getRandomValueFromInterval(b.opts.randomizationFactor, currentInterval)

			time.Sleep(currentInterval)
		} else {
			return nil
		}
	}
}
func getRandomValueFromInterval(randomizationFactor float64, currentInterval time.Duration) time.Duration {
	delta := randomizationFactor * float64(currentInterval)
	minInterval := float64(currentInterval) - delta
	maxInterval := float64(currentInterval) + delta
	return time.Duration(minInterval + (rand.New(rand.NewSource(time.Now().UnixNano())).Float64() * (maxInterval - minInterval + 1)))
}
