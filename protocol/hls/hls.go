package hls

import (
	"fmt"
	"net/http"
	"strconv"
)

const (
	duration = 3000
)

var KeepAfterEnd = true

var (
	ErrNoPublisher         = fmt.Errorf("no publisher")
	ErrInvalidReq          = fmt.Errorf("invalid req url path")
	ErrNoSupportVideoCodec = fmt.Errorf("no support video codec")
	ErrNoSupportAudioCodec = fmt.Errorf("no support audio codec")
)

func HandleM3u8Writer(key string, getSource func(string) (*Source, bool), w http.ResponseWriter) {
	source, ok := getSource(key)
	if !ok {
		http.Error(w, ErrNoPublisher.Error(), http.StatusForbidden)
		return
	}
	tsCache := source.GetCacheInc()
	if tsCache == nil {
		http.Error(w, ErrNoPublisher.Error(), http.StatusForbidden)
		return
	}
	body, err := tsCache.GenM3U8PlayList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "application/x-mpegURL")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.Write(body)
}

func HandleTsWriter(key, url string, getSource func(string) (*Source, bool), w http.ResponseWriter) {
	source, ok := getSource(key)
	if !ok {
		http.Error(w, ErrNoPublisher.Error(), http.StatusForbidden)
		return
	}
	tsCache := source.GetCacheInc()
	item, err := tsCache.GetItem(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "video/mp2ts")
	w.Header().Set("Content-Length", strconv.Itoa(len(item.Data)))
	w.Write(item.Data)
}
