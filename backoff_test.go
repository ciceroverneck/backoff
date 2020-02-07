package backoff

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestBackoff_Do(t *testing.T) {
	ctx := context.Background()
	type args struct {
		fn Handle
	}
	spectedErr := errors.New("teste")
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test Continue",
			args: args{
				fn: func(ctx context.Context) error {
					fmt.Println("teste")
					return Continue(spectedErr)
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := New(Exponential(), MaxRetries(3))
			if err := this.Do(ctx, tt.args.fn); (err != nil) && err != spectedErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
