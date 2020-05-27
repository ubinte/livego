package main

import (
	"testing"

	"github.com/ubinte/livego/app"
	"github.com/ubinte/livego/protocol/rtmp"
)

func TestStartHttpflvServer(t *testing.T) {
	app.AddApp("live").AddChannelKey("insecure_channel_key", "movie")
	stream := rtmp.NewRtmpStream()
	go StartRtmpServer(stream, ":1935") // push rtmp://127.0.0.1:1935/live/insecure_channel_key
	StartHttpflvServer(stream, ":80")   // pull http://127.0.0.1/live/movie.flv
}
