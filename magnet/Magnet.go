package magnet

import (
	"fmt"
	"net/url"
	"strings"
)

type Magnet struct {
	InfoHash    string
	DisplayName string
	Trackers    []string
	WebSeeds    []string
	ExactSource string
	Extra       map[string][]string
}

func ParseMagnet(link string) (*Magnet, error) {
	m := Magnet{
		Extra: make(map[string][]string),
	}
	query := strings.Split(link, "?")

	if len(query) != 2 {
		return nil, fmt.Errorf("invalid info magnet link")
	}

	params := strings.Split(query[1], "&")

	for _, s := range params {
		kv := strings.Split(s, "=")

		key := kv[0]

		value, err := url.QueryUnescape(kv[1])
		if err != nil {
			value = kv[1]
		}

		switch key {
		case "xt":
			if strings.HasPrefix(value, "urn:btih:") {
				m.InfoHash = strings.TrimPrefix(value, "urn:btih:")
			}

		case "dn":
			m.DisplayName = value
		case "tr":
			m.Trackers = append(m.Trackers, value)
		case "ws":
			m.WebSeeds = append(m.WebSeeds, value)
		case "xs":
			m.ExactSource = value
		default:
			m.Extra[key] = append(m.Extra[key], value)

		}
	}
	return &m, nil
}
