package httpflv

import (
	"net/http"

	"github.com/ubinte/livego/av"
)

func HandleWriter(app, channel, url string, handler av.Handler, w http.ResponseWriter) {
	flvWriter := NewFLVWriter(app, channel, url, w)
	handler.HandleWriter(flvWriter)
	flvWriter.Wait()
}
