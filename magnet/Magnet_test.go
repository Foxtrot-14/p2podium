package magnet

import (
	"reflect"
	"testing"
)

func TestParseMagnet(t *testing.T) {
	tests := map[string]struct {
		link   string
		want   *Magnet
		hasErr bool
	}{
		"empty link": {
			link:   "",
			want:   nil,
			hasErr: true,
		},
		"simple magnet": {
			link: "magnet:?xt=urn:btih:12345&dn=testfile",
			want: &Magnet{
				InfoHash:    "12345",
				DisplayName: "testfile",
				Trackers:    nil,
				WebSeeds:    nil,
				ExactSource: "",
				Extra:       make(map[string][]string),
			},
			hasErr: false,
		},
		"magnet with trackers and webseeds": {
			link: "magnet:?xt=urn:btih:abcdef&dn=myfile&tr=udp://tracker1&tr=udp://tracker2&ws=http://webseed",
			want: &Magnet{
				InfoHash:    "abcdef",
				DisplayName: "myfile",
				Trackers:    []string{"udp://tracker1", "udp://tracker2"},
				WebSeeds:    []string{"http://webseed"},
				ExactSource: "",
				Extra:       make(map[string][]string),
			},
			hasErr: false,
		},
		"magnet with extra params": {
			link: "magnet:?xt=urn:btih:xyz&kt=tag1&kt=tag2&as=http://example",
			want: &Magnet{
				InfoHash:    "xyz",
				Trackers:    nil,
				WebSeeds:    nil,
				ExactSource: "",
				Extra: map[string][]string{
					"kt": {"tag1", "tag2"},
					"as": {"http://example"},
				},
			},
			hasErr: false,
		},
		"invalid magnet": {
			link:   "not-a-magnet-link",
			want:   nil,
			hasErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ParseMagnet(tt.link)
			if (err != nil) != tt.hasErr {
				t.Fatalf("expected error: %v, got: %v", tt.hasErr, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("expected %+v, got %+v", tt.want, got)
			}
		})
	}
}


