package main

import (
	"testing"

	"github.com/ubinte/livego/app"
	"github.com/ubinte/livego/protocol/rtmp"
)

func TestStartRtmpClient(t *testing.T) {
	app.AddApp("live").AddChannelKey("insecure_channel_key", "movie")
	stream := rtmp.NewRtmpStream()
	go StartRtmpServer(stream, ":1935")                                         // push rtmp://127.0.0.1:1935/live/insecure_channel_key
	StartRtmpClient(stream, "rtmp://127.0.0.1:1935/live/moive", "tmp_flvCache") // pull to /flvCache
}
