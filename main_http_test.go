package main

import (
	"net/http"
	"path"
	"testing"

	"github.com/ubinte/livego/app"
	"github.com/ubinte/livego/protocol/hls"
	"github.com/ubinte/livego/protocol/httpflv"
	"github.com/ubinte/livego/protocol/rtmp"
)

func TestStartHttpServer(t *testing.T) {
	app.AddApp("live").AddChannelKey("insecure_channel_key", "movie")

	stream := rtmp.NewRtmpStream()
	go StartRtmpServer(stream, ":1935")

	hlsServer := hls.NewServer(stream)

	http.HandleFunc("/flv/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		httpflv.HandleWriter("live", "movie", r.URL.Path, stream, w)
	})
	http.HandleFunc("/hls/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if path.Ext(r.URL.Path) == ".m3u8" {
			hls.HandleM3u8Writer("live/movie", func(key string) (*hls.Source, bool) {
				s, ok := hlsServer.GetSource(key)
				if ok {
					s.TsFilePath = "/hls/"
					return s, true
				} else {
					return nil, false
				}
			}, w)
		} else if path.Ext(r.URL.Path) == ".ts" {
			hls.HandleTsWriter("live/movie", r.URL.Path, hlsServer.GetSource, w)
		}
	})
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		key := app.AddApp("live").AddChannel("carton")
		w.Write([]byte(key))
	})
	http.ListenAndServe(":80", nil)
}
