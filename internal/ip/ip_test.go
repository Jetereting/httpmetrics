package ip

import (
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_loadFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name     string
		args     args
		wantList *AllowList
		wantErr  bool
	}{
		{
			name: "readFile",
			args: args{filename: "./testdata/allowip.toml"},
			wantList: &AllowList{
				IPs: []net.IP{
					net.IPv4(192, 168, 1, 200),
				},
				IPNets: []*net.IPNet{
					{
						IP:   net.IPv4(127, 0, 0, 0),
						Mask: net.IPv4Mask(255, 0, 0, 0),
					},
					{
						IP:   net.IPv4(192, 168, 1, 0),
						Mask: net.IPv4Mask(255, 255, 255, 240),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				filename: "nil",
			},
			wantList: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotList, err := loadFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantList != nil {
				tt.wantList.sort()
			}
			if !cmp.Equal(gotList, tt.wantList) {
				t.Errorf("loadFile() = %v, want %v", gotList, tt.wantList)
			}
		})
	}
}

func Test_loadBytes(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name     string
		args     args
		wantList *AllowList
		wantErr  bool
	}{
		{
			name: "readBytes",
			args: args{
				b: []byte(`ips = [
    '127.0.0.1',
    '127.0.0.1/8',
    '127.0.0.1/24',
    '192.168.1.200',
    '192.168.1.2/32',
    '172.32.1.2',
	'192.168.1.6/28',
	'192',
	'::1',
	'fec0::b:20c:29ff:fe0b:55c3'
]`),
			},
			wantList: &AllowList{
				IPs: []net.IP{
					net.ParseIP("192.168.1.200"),
					net.IPv6loopback,
					net.IPv4(172, 32, 1, 2),
					net.ParseIP("fec0::b:20c:29ff:fe0b:55c3"),
				},
				IPNets: []*net.IPNet{
					{
						IP:   net.IPv4(127, 0, 0, 0),
						Mask: net.IPv4Mask(255, 0, 0, 0),
					},
					{
						IP:   net.IPv4(192, 168, 1, 0),
						Mask: net.IPv4Mask(255, 255, 255, 240),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalidBytes",
			args: args{
				b: []byte(`]`),
			},
			wantList: nil,
			wantErr:  true,
		},
		{
			name: "nilBytes",
			args: args{
				b: nil,
			},
			wantList: nil,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotList, err := loadBytes(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantList != nil {
				tt.wantList.sort()
			}
			if !cmp.Equal(gotList, tt.wantList) {
				t.Errorf("loadBytes() = %v, want %v", gotList, tt.wantList)
			}
		})
	}
}

func Test_loadArray(t *testing.T) {
	type args struct {
		a []string
	}
	tests := []struct {
		name string
		args args
		want *AllowList
	}{
		{
			name: "readArray",
			args: args{
				a: []string{
					"127.0.0.1",
					"127.0.0.1/8",
					"127.0.0.1/24",
					"192.168.1.200",
					"192.168.1.2/32",
					"192.168.1.6/28",
					"::1",
				},
			},
			want: &AllowList{
				IPs: []net.IP{
					net.IPv4(192, 168, 1, 200),
					net.IPv6loopback,
				},
				IPNets: []*net.IPNet{
					{
						IP:   net.IPv4(127, 0, 0, 0),
						Mask: net.IPv4Mask(255, 0, 0, 0),
					},
					{
						IP:   net.IPv4(192, 168, 1, 0),
						Mask: net.IPv4Mask(255, 255, 255, 240),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != nil {
				tt.want.sort()
			}
			if got := LoadArray(tt.args.a); !cmp.Equal(got, tt.want) {
				t.Errorf("loadArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList_ContainsString(t *testing.T) {
	type fields struct {
		IPs    []net.IP
		IPNets []*net.IPNet
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "contains",
			fields: fields{
				IPs: []net.IP{
					net.IPv4(192, 168, 1, 200),
					net.IPv6loopback,
				},
				IPNets: []*net.IPNet{
					{
						IP:   net.IPv4(127, 0, 0, 0),
						Mask: net.IPv4Mask(255, 0, 0, 0),
					},
					{
						IP:   net.IPv4(192, 168, 1, 0),
						Mask: net.IPv4Mask(255, 255, 255, 240),
					},
				},
			},
			args: args{
				s: "192.168.1.200",
			},
			want: true,
		},
		{
			name: "notContains",
			fields: fields{
				IPs: []net.IP{
					net.IPv4(192, 168, 1, 200),
					net.IPv6loopback,
				},
				IPNets: []*net.IPNet{
					{
						IP:   net.IPv4(127, 0, 0, 0),
						Mask: net.IPv4Mask(255, 0, 0, 0),
					},
					{
						IP:   net.IPv4(192, 168, 1, 0),
						Mask: net.IPv4Mask(255, 255, 255, 240),
					},
				},
			},
			args: args{
				s: "192.168.1",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &AllowList{
				IPs:    tt.fields.IPs,
				IPNets: tt.fields.IPNets,
			}
			if got := l.ContainsString(tt.args.s); got != tt.want {
				t.Errorf("List.ContainsString() = %v, want %v", got, tt.want)
			}
		})
	}
}
