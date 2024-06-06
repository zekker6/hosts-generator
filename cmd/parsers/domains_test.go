package parsers

import (
	"reflect"
	"testing"
)

func TestFilterWildcardEntries(t *testing.T) {
	type args struct {
		domains []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Wildcard test",
			args: args{
				domains: []string{
					"*.test.com",
					"test.com",
					"*.test.com",
					"test.com",
				},
			},
			want: []string{
				"test.com",
				"test.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterWildcardEntries(tt.args.domains); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterWildcardEntries() = %v, want %v", got, tt.want)
			}
		})
	}
}
