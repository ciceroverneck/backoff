package backoff

import (
	"reflect"
	"testing"
	"time"

	"ciceroverneck.dev/backoff/internal"
)

func TestExponential(t *testing.T) {
	tests := []struct {
		opt  Option
		want bool
	}{
		{
			opt:  Exponential(),
			want: true,
		},
	}

	for _, test := range tests {
		got := internal.BackoffSettings{}
		test.opt(&got)
		if got.Exponential != test.want {
			t.Errorf("BackoffSettings.Exponential = %v, want %v", got.Exponential, test.want)
		}
	}
}

func TestInterval(t *testing.T) {
	tests := []struct {
		opt  Option
		want time.Duration
	}{
		{
			opt:  Interval(time.Second * 3),
			want: time.Second * 3,
		},
		{
			opt:  Interval(time.Second),
			want: time.Second,
		},
	}

	for _, test := range tests {
		got := internal.BackoffSettings{}
		test.opt(&got)
		if got.Interval != test.want {
			t.Errorf("BackoffSettings.Interval = %v, want %v", got.Interval, test.want)
		}
	}
}

func TestMaxElapsedTime(t *testing.T) {
	tests := []struct {
		opt  Option
		want time.Duration
	}{
		{
			opt:  MaxElapsedTime(time.Second * 3),
			want: time.Second * 3,
		},
		{
			opt:  MaxElapsedTime(time.Second),
			want: time.Second,
		},
	}

	for _, test := range tests {
		got := internal.BackoffSettings{}
		test.opt(&got)
		if got.MaxElapsedTime != test.want {
			t.Errorf("BackoffSettings.MaxElapsedTime = %v, want %v", got.MaxElapsedTime, test.want)
		}
	}
}

func TestMaxInterval(t *testing.T) {
	tests := []struct {
		opt  Option
		want time.Duration
	}{
		{
			opt:  MaxInterval(time.Second * 3),
			want: time.Second * 3,
		},
		{
			opt:  MaxInterval(time.Second),
			want: time.Second,
		},
	}

	for _, test := range tests {
		got := internal.BackoffSettings{}
		test.opt(&got)
		if got.MaxInterval != test.want {
			t.Errorf("BackoffSettings.MaxInterval = %v, want %v", got.MaxInterval, test.want)
		}
	}
}

func TestMaxRetries(t *testing.T) {
	tests := []struct {
		opt  Option
		want uint
	}{
		{
			opt:  MaxRetries(10),
			want: 10,
		},
		{
			opt:  MaxRetries(2220),
			want: 2220,
		},
	}

	for _, test := range tests {
		got := internal.BackoffSettings{}
		test.opt(&got)
		if got.MaxRetries != test.want {
			t.Errorf("BackoffSettings.MaxRetries = %v, want %v", got.MaxRetries, test.want)
		}
	}
}

func TestMultiplier(t *testing.T) {
	tests := []struct {
		opt  Option
		want float64
	}{
		{
			opt:  Multiplier(10.0),
			want: 10.0,
		},
		{
			opt:  Multiplier(2.220),
			want: 2.220,
		},
	}

	for _, test := range tests {
		got := internal.BackoffSettings{}
		test.opt(&got)
		if got.Multiplier != test.want {
			t.Errorf("BackoffSettings.Multiplier = %v, want %v", got.Multiplier, test.want)
		}
	}
}

func TestNotify(t *testing.T) {
	fn_cast := NotifyHandle(func(err error, duration time.Duration, count uint) {})
	fn := func(err error, duration time.Duration, count uint) {}
	tests := []struct {
		opt  Option
		want NotifyHandle
	}{
		{
			opt:  Notify(fn),
			want: fn,
		},
		{
			opt:  Notify(fn_cast),
			want: fn_cast,
		},
		{
			opt:  Notify(nil),
			want: nil,
		},
	}

	for _, test := range tests {
		got := internal.BackoffSettings{}
		test.opt(&got)
		if reflect.ValueOf(got.Notify).Pointer() != reflect.ValueOf(test.want).Pointer() {
			t.Errorf("BackoffSettings.Notify = %p, want %p", got.Notify, test.want)
		}
	}
}

func TestRandomizationFactor(t *testing.T) {
	tests := []struct {
		opt  Option
		want float64
	}{
		{
			opt:  RandomizationFactor(22.22),
			want: 22.22,
		},
		{
			opt:  RandomizationFactor(1),
			want: 1,
		},
		{
			opt:  RandomizationFactor(0.5),
			want: 0.5,
		},
	}

	for _, test := range tests {
		got := internal.BackoffSettings{}
		test.opt(&got)
		if got.RandomizationFactor != test.want {
			t.Errorf("BackoffSettings.RandomizationFactor = %v, want %v", got.RandomizationFactor, test.want)
		}
	}
}
