package base58

import (
	"testing"
)

func TestEncodeFromInt(t *testing.T) {
	type args struct {
		i int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Encode 1000000000",
			args: args{
				i: 1000000000,
			},
			want: "2XNGAK",
		},
		{
			name: "Encode 1008360000",
			args: args{
				i: 1008360000,
			},
			want: "2Y77JF",
		},
		{
			name: "Encode 32000930910",
			args: args{
				i: 32000930910,
			},
			want: "qkp8s7",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeFromInt(tt.args.i); got != tt.want {
				t.Errorf("EncodeFromInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeToInt(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Decode 2Y77JF",
			args: args{s: "2Y77JF"},
			want: 1008360000,
			wantErr: false,
		},
		{
			name: "Decode 2XNGAK",
			args: args{s: "2XNGAK"},
			want: 1000000000,
			wantErr: false,
		},
		{
			name: "Decode qkp8s7",
			args: args{s: "qkp8s7"},
			want: 32000930910,
			wantErr: false,
		},
		{
			name: "Decode with incorrect string",
			args: args{s: "Omn38u"},
			want: 0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeToInt(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if got != tt.want {
					t.Errorf("DecodeToInt() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
