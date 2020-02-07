package backoff

import (
	"time"

	"ciceroverneck.dev/backoff/internal"
)


type Option func(*internal.BackoffSettings)

func MaxRetries(v uint) Option {
	return func(args *internal.BackoffSettings) {
		args.MaxRetries = v
	}
}

func Exponential() Option {
	return func(args *internal.BackoffSettings) {
		args.Exponential = true
	}
}

func Notify(fn NotifyHandle) Option {
	return func(args *internal.BackoffSettings) {
		args.Notify = fn
	}
}

func MaxElapsedTime(duration time.Duration) Option {
	return func(args *internal.BackoffSettings) {
		args.MaxElapsedTime = duration
	}
}

func RandomizationFactor(randomizationFactor float64) Option {
	return func(args *internal.BackoffSettings) {
		args.RandomizationFactor = randomizationFactor
	}
}

func Multiplier(multiplier float64) Option {
	return func(args *internal.BackoffSettings) {
		args.Multiplier = multiplier
	}
}

func MaxInterval(duration time.Duration) Option {
	return func(args *internal.BackoffSettings) {
		args.MaxInterval = duration
	}
}

func Interval(duration time.Duration) Option {
	return func(args *internal.BackoffSettings) {
		args.Interval = duration
	}
}