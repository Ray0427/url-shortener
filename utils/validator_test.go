package utils

import "testing"

func TestCheckUrl(t *testing.T) {
	type args struct {
		fullUrl string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "ValidUrl",
			args: args{
				fullUrl: "https://ipinfo.io",
			},
			want: true,
		},
		{
			name: "InvalidUrl",
			args: args{
				fullUrl: "ipinfo.io",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckUrl(tt.args.fullUrl); got != tt.want {
				t.Errorf("CheckUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
