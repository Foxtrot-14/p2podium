package magnet

import (
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
)

type Magnet struct {
	InfoHash    [20]byte
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

	query := strings.SplitN(link, "?", 2)
	if len(query) != 2 {
		return nil, fmt.Errorf("invalid magnet link: missing '?'")
	}

	params := strings.Split(query[1], "&")
	for _, s := range params {
		kv := strings.SplitN(s, "=", 2)
		if len(kv) != 2 {
			continue
		}

		key := kv[0]
		value, err := url.QueryUnescape(kv[1])
		if err != nil {
			value = kv[1]
		}

		switch key {
		case "xt":
			if strings.HasPrefix(value, "urn:btih:") {
				hashStr := strings.TrimPrefix(value, "urn:btih:")

				var infohash []byte
				if len(hashStr) == 40 {
					infohash, err = hex.DecodeString(hashStr)
				} else if len(hashStr) == 32 {
					infohash, err = base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(hashStr))
				} else {
					err = fmt.Errorf("invalid infohash length: %d", len(hashStr))
				}

				if err != nil {
					return nil, fmt.Errorf("failed to decode infohash: %v", err)
				}

				copy(m.InfoHash[:], infohash)
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
