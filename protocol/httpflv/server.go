package httpflv

import (
	"net"
	"net/http"
	"path"
	"strings"

	"github.com/ubinte/livego/protocol/rtmp"
)

type Server struct {
	handler *rtmp.RtmpStream
}

func NewServer(h *rtmp.RtmpStream) *Server {
	return &Server{
		handler: h,
	}
}

func (self *Server) Serve(l net.Listener) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", self.HandleConn)
	http.Serve(l, mux)
	return nil
}

func (self *Server) HandleConn(w http.ResponseWriter, r *http.Request) {
	if path.Ext(r.URL.Path) != ".flv" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	appName := strings.TrimPrefix(path.Dir(r.URL.Path), "/")
	channelName := strings.TrimSuffix(path.Base(r.URL.Path), path.Ext(r.URL.Path))
	key := appName + "/" + channelName

	publishers := self.handler.GetReaders()
	if _, ok := publishers[key]; !ok {
		http.Error(w, "channel is close", http.StatusNotFound)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	HandleWriter(appName, channelName, r.URL.Path, self.handler, w)
}
