package rtmp

import (
	"fmt"
	"net"

	"github.com/ubinte/livego/app"
	"github.com/ubinte/livego/av"
	"github.com/ubinte/livego/protocol/rtmp/core"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	handler     av.Handler
	AllowClient bool
}

func NewServer(h av.Handler) *Server {
	return &Server{
		handler:     h,
		AllowClient: false,
	}
}

func (self *Server) Serve(listener net.Listener) error {
	defer func() {
		if r := recover(); r != nil {
			log.Error("rtmp serve panic: ", r)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		rmtpConn := core.NewConn(conn, 4*1024)
		log.Debug("new client, connect remote: ", rmtpConn.RemoteAddr().String(),
			"local:", rmtpConn.LocalAddr().String())
		go self.handleConn(rmtpConn)
	}
}

func (self *Server) handleConn(conn *core.Conn) error {
	if err := conn.HandshakeServer(); err != nil {
		conn.Close()
		log.Error("handleConn HandshakeServer err: ", err)
		return err
	}
	connServer := core.NewConnServer(conn)

	if err := connServer.ReadMsg(); err != nil {
		conn.Close()
		log.Error("handleConn read msg err: ", err)
		return err
	}

	appName, channelKey, _ := connServer.GetInfo()
	myApp, ok := app.GetApp(appName)
	if !ok {
		conn.Close()
		err := fmt.Errorf("invalid app:", appName)
		log.Error(err)
		return err
	}

	log.Debugf("handleConn: IsPublisher=%v", connServer.IsPublisher())
	if connServer.IsPublisher() {
		channel, ok := myApp.GetChannel(channelKey)
		if !ok {
			conn.Close()
			err := fmt.Errorf("invalid channel key:", channelKey)
			log.Error(err.Error())
			return err
		}
		connServer.PublishInfo.Name = channel

		reader := NewVirReader(connServer)
		self.handler.HandleReader(reader)
		log.Debugf("new publisher: %+v", reader.Info())

	} else if self.AllowClient {
		writer := NewVirWriter(connServer)
		log.Debugf("new player: %+v", writer.Info())
		self.handler.HandleWriter(writer)
	}

	return nil
}
