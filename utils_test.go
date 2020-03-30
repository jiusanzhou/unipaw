package unipaw

import "testing"

func TestPathJoin(t *testing.T) {
	type args struct {
		a []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"normal 1", args{[]string{"a", "b", "c"}}, "a/b/c"},
		{"normal 2", args{[]string{""}}, "/"},
	}
	for _, tt := range tests {
		if got := PathJoin(tt.args.a...); got != tt.want {
			t.Errorf("%q. PathJoin() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
