package hls

import (
	"fmt"
	"net"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/ubinte/livego/av"
	"github.com/ubinte/livego/protocol/rtmp"

	cmap "github.com/orcaman/concurrent-map"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	handler *rtmp.RtmpStream
	sources cmap.ConcurrentMap
}

func NewServer(h *rtmp.RtmpStream) *Server {
	server := &Server{
		handler: h,
		sources: cmap.New(),
	}
	h.AddGetWriter(server)

	go server.checkStop()
	return server
}

func (server *Server) Serve(listener net.Listener) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server.Handle(w, r)
	})
	http.Serve(listener, mux)
	return nil
}

func (self *Server) GetWriter(info av.Info) av.WriteCloser {
	v, ok := self.sources.Get(info.Key)
	if ok {
		return v.(*Source)
	} else {
		return self.NewSource(info)
	}
}

func (self *Server) NewSource(info av.Info) *Source {
	source := NewSource(info)
	self.sources.Set(info.Key, source)
	return source
}

func (self *Server) GetSource(key string) (*Source, bool) {
	v, ok := self.sources.Get(key)
	if ok {
		return v.(*Source), true
	} else {
		return nil, false
	}
}

func (self *Server) checkStop() {
	for {
		<-time.After(5 * time.Second)
		for item := range self.sources.IterBuffered() {
			v := item.Val.(*Source)
			if !v.Alive() && !KeepAfterEnd {
				log.Debug("check stop and remove: ", v.Info())
				self.sources.Remove(item.Key)
			}
		}
	}
}

func (self *Server) HandleCrossdomain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	crossdomainxml := []byte(`<?xml version="1.0" ?><cross-domain-policy><allow-access-from domain="*" /><allow-http-request-headers-from domain="*" headers="*"/></cross-domain-policy>`)
	w.Write(crossdomainxml)
}

func (self *Server) Handle(w http.ResponseWriter, r *http.Request) {
	if path.Base(r.URL.Path) == "crossdomain.xml" {
		self.HandleCrossdomain(w, r)
		return
	}

	switch path.Ext(r.URL.Path) {
	case ".m3u8":
		key, _ := self.parseM3u8(r.URL.Path)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		HandleM3u8Writer(key, self.GetSource, w)
	case ".ts":
		key, _ := self.parseTs(r.URL.Path)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		HandleTsWriter(key, r.URL.Path, self.GetSource, w)
	}
}

func (server *Server) parseM3u8(pathstr string) (key string, err error) {
	pathstr = strings.TrimLeft(pathstr, "/")
	key = strings.Split(pathstr, path.Ext(pathstr))[0]
	return
}

func (server *Server) parseTs(pathstr string) (key string, err error) {
	pathstr = strings.TrimLeft(pathstr, "/")
	paths := strings.SplitN(pathstr, "/", 3)
	if len(paths) != 3 {
		err = fmt.Errorf("invalid path=%s", pathstr)
		return
	}
	key = paths[0] + "/" + paths[1]

	return
}
