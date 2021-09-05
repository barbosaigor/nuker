package pipeline

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalcRequests(t *testing.T) {
	type args struct {
		start time.Duration
		end   time.Duration
		curr  time.Duration
		min   int
		max   int
	}

	type test struct {
		name string
		args args
		want int
	}

	tests := []test{
		{
			name: "when empty args, should return zero",
			args: args{
				start: 0,
				end:   0,
				curr:  0,
				min:   0,
				max:   0,
			},
			want: 0,
		},
		{
			name: "when have reached around 30%, should return proportional value",
			args: args{
				start: 1 * time.Second,
				end:   3 * time.Second,
				curr:  1100 * time.Millisecond,
				min:   1,
				max:   10,
			},
			want: 1,
		},
		{
			name: "when in half way, should return half of max",
			args: args{
				start: 1 * time.Second,
				end:   3 * time.Second,
				curr:  2 * time.Second,
				min:   1,
				max:   10,
			},
			want: 5,
		},
		{
			name: "when have reached around 80%, should return proportional value",
			args: args{
				start: 1 * time.Second,
				end:   3 * time.Second,
				curr:  2600 * time.Millisecond,
				min:   1,
				max:   10,
			},
			want: 8,
		},
		{
			name: "when completed, should return max",
			args: args{
				start: 1 * time.Second,
				end:   3 * time.Second,
				curr:  3 * time.Second,
				min:   1,
				max:   10,
			},
			want: 10,
		},
		{
			name: "when max is less than min, should use min as max",
			args: args{
				start: 1 * time.Second,
				end:   3 * time.Second,
				curr:  2 * time.Second,
				min:   1,
				max:   -10,
			},
			want: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := &pipeline{}

			got := p.calcRequests(
				tc.args.start,
				tc.args.end,
				tc.args.curr,
				tc.args.min,
				tc.args.max,
			)

			assert.Equal(t, tc.want, got, "calcRequests(...) = %v, want = %v", got, tc.want)
		})
	}
}
