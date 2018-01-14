package hermes

import (
	"reflect"
	"testing"
	"time"
)

func Test_newOptions(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want Options
	}{
		{
			name: "Default",
			want: Options{
				url: defaultSocketURL,
			},
		},
		{
			name: "With deadlock",
			args: args{
				opts: []Option{
					Deadline(time.Second),
				},
			},
			want: Options{
				url:      defaultSocketURL,
				deadline: time.Second,
			},
		},
		{
			name: "With deadlock and url",
			args: args{
				opts: []Option{
					URL("test"),
					Deadline(time.Second),
				},
			},
			want: Options{
				url:      "test",
				deadline: time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newOptions(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
