package magnet

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/samsterlicious/torrent/tracker"
)

type Magnet struct {
	FileName string
	Hash     string
	trackers []string
}

func ParseMagnetLink(link string) (*Magnet, error) {
	r := regexp.MustCompile("xt=urn:btih:([^&]+)")
	hashSlice := r.FindStringSubmatch(link)
	if len(hashSlice) < 2 {
		return nil, errors.New("no hash present in magnet link")
	}
	r = regexp.MustCompile(`tr=([^&]+)`)
	allMatches := r.FindAllStringSubmatch(link, -1)
	trackers := make([]string, len(allMatches))
	for i, match := range allMatches {
		decodedValue, _ := url.QueryUnescape(match[1])
		trackers[i] = decodedValue
	}

	r = regexp.MustCompile("dn=([^&]+)")
	fileName := r.FindStringSubmatch(link)

	if len(fileName) < 2 {
		return nil, errors.New("no file name in magnet link")
	}

	decodedValue, _ := url.QueryUnescape(fileName[1])

	ret := Magnet{decodedValue, hashSlice[1], trackers}
	return &ret, nil
}

func (mag *Magnet) SendConnectionRequest() {
	responseChan := make(chan []byte)

	for _, link := range mag.trackers {
		if tracker.IsUdp(link) {
			go tracker.ProcessUdp(link, responseChan)
			break
		}
	}

	resp := <-responseChan
	fmt.Println(resp)
}
