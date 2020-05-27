package app

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var apps = make(map[string]application)

func AddApp(name string) *application {
	a := application{
		Static:   make([]string, 0),
		channels: cache.New(cache.DefaultExpiration, 30*time.Second),
	}
	apps[name] = a
	return &a
}

func GetApp(name string) (application, bool) {
	app, ok := apps[name]
	return app, ok
}

func Find(appName, channelName string) (string, bool) {
	app, ok := GetApp(appName)
	if !ok {
		return "", false
	}

	channelKey, ok := app.GetChannel(channelName)
	if !ok {
		return "", false
	}

	return channelKey, true
}
