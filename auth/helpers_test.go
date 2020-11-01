package auth

import (
	"fmt"
	"testing"
)

func Test_randString(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name       string
		args       args
		wantLength int
	}{
		{
			name: "positive_rand_string",
			args: args{
				n: 20,
			},
			wantLength: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := randString(tt.args.n)
			fmt.Println(got)
			if len(got) != tt.wantLength {
				t.Errorf("randString() length = %d (%v), want %d", len(got), got, tt.wantLength)
			}
		})
	}
}
