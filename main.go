package main

import (
	"net"
	"time"

	"github.com/ubinte/livego/app"
	"github.com/ubinte/livego/av"
	"github.com/ubinte/livego/container/flv"
	"github.com/ubinte/livego/protocol/hls"
	"github.com/ubinte/livego/protocol/httpflv"
	"github.com/ubinte/livego/protocol/rtmp"

	log "github.com/sirupsen/logrus"
)

var VERSION = "0.1"
var CONFIG = "livego.json"

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	app.AddApp("live").AddChannelKey("insecure_channel_key", "movie")

	stream := rtmp.NewRtmpStream()
	go StartRtmpServer(stream, ":1935")
	StartHttpflvServer(stream, ":80")
}

func StartHlsServer(stream *rtmp.RtmpStream, addr string) {
	hlsServer := hls.NewServer(stream)
	hlsListener, _ := net.Listen("tcp", addr)
	hlsServer.Serve(hlsListener)
}

func StartHttpflvServer(stream *rtmp.RtmpStream, addr string) {
	httpflvServer := httpflv.NewServer(stream)
	httpflvListener, _ := net.Listen("tcp", addr)
	httpflvServer.Serve(httpflvListener)
}

func StartRtmpServer(stream *rtmp.RtmpStream, addr string) {
	rtmpServer := rtmp.NewServer(stream)
	rmtpListener, _ := net.Listen("tcp", addr)
	rtmpServer.Serve(rmtpListener)
}

func StartRtmpClient(stream *rtmp.RtmpStream, addr, path string) {
	rmtpClient := rtmp.NewClient(stream, flv.NewFlvDvr(path))
	if err := rmtpClient.Dial(addr, av.PLAY); err != nil {
		log.Error(err)
	}

	for {
		<-time.After(100 * time.Millisecond)
	}
}
