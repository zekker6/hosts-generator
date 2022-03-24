package traefik_v2

import (
	"reflect"
	"testing"
)

func Test_extractHosts(t *testing.T) {
	type args struct {
		rules []string
	}

	tc := NewTraefikV2Client("localhost")
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Rules test",
			args: args{
				rules: []string{
					"Host(`gateway.gateway.docker.localhost`)",
					"Host(`test.com.asd`)",
					"Host(`Some.WWW.test.com.asd`)",
					"Host(`test`);PathPrefix(`/api/`)",
					"Host(`test`) && PathPrefix:/api/",
					"Host(`test.some.asdf`) || Host(`test.com.zxc`) || Host(`some.test.zxc`)",
				},
			},
			want: []string{
				"Some.WWW.test.com.asd",
				"gateway.gateway.docker.localhost",
				"some.test.zxc",
				"test",
				"test",
				"test.com.asd",
				"test.com.zxc",
				"test.some.asdf",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tc.extractHosts(tt.args.rules); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractHosts() = %v, want %v", got, tt.want)
			}
		})
	}
}
