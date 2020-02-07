package internal

import "time"

type BackoffSettings struct {
	MaxRetries          uint
	Exponential         bool
	Notify              func(err error, duration time.Duration, count uint)
	Multiplier          float64
	RandomizationFactor float64
	MaxElapsedTime      time.Duration
	Interval            time.Duration
	MaxInterval         time.Duration
}
