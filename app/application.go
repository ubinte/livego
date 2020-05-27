package app

import (
	"github.com/patrickmn/go-cache"
	"github.com/ubinte/goutils/strutils"
)

type application struct {
	Static   []string
	channels *cache.Cache
}

func (self *application) AddChannelKey(key, name string) {
	self.channels.Add(key, name, cache.DefaultExpiration)
}

func (self *application) AddChannel(name string) string {
	key := strutils.RandomANI(32)
	self.channels.Add(key, name, cache.DefaultExpiration)
	return key
}

func (self *application) GetChannel(key string) (string, bool) {
	if name, ok := self.channels.Get(key); ok {
		return name.(string), true
	} else {
		return "", false
	}
}
