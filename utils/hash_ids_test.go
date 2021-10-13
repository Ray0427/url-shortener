package utils

import "testing"

func TestEncode(t *testing.T) {
	type args struct {
		salt      string
		minLength int
		num       int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "TestEncodeSuccess",
			args: args{
				salt:      "test",
				minLength: 10,
				num:       1,
			},
			want: "3wedgpzLRq",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Encode(tt.args.salt, tt.args.minLength, tt.args.num); got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	type args struct {
		salt      string
		minLength int
		hash      string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestDecodeSuccess",
			args: args{
				salt:      "test",
				minLength: 10,
				hash:      "3wedgpzLRq",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "TestDecodeFailed",
			args: args{
				salt:      "test",
				minLength: 10,
				hash:      "3wedgpzLR",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.args.salt, tt.args.minLength, tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}
