package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/gogap/cache_storages"
	"github.com/gogap/logs"
	"github.com/gogap/spirit"

	"github.com/spirit-contrib/session/controllers"
)

type SessionConfig struct {
	MemcachedConfig `json:"memcached"`
}

type MemcachedConfig struct {
	Address string `json:"address"`
}

func main() {
	logs.SetFileLogger("./logs/session.log")

	storageAddr := ""
	if bFile, e := ioutil.ReadFile("./conf/session.conf"); e != nil {
		panic(e)
	} else {
		sessionConfig := SessionConfig{}
		if e := json.Unmarshal(bFile, &sessionConfig); e != nil {
			panic(e)
		}

		sessionConfig.Address = strings.TrimSpace(sessionConfig.Address)
		if sessionConfig.Address == "" {
			panic("memcached.address is empty")
		}

		storageAddr = sessionConfig.Address
	}

	var sessionStorage *controllers.SessionStorage
	if storage, e := cache_storages.NewMemcachedStorage(storageAddr); e != nil {
		panic(e)
	} else {
		sessionStorage = controllers.NewSessionStorage(storage)
	}

	sessionSpirit := spirit.NewClassicSpirit("session", "a basic session component", "1.0.0")
	sessionComponent := spirit.NewBaseComponent("session")

	sessionComponent.RegisterHandler("set_session", sessionStorage.SetSession)
	sessionComponent.RegisterHandler("get_session", sessionStorage.GetSession)

	sessionSpirit.Hosting(sessionComponent).Build().Run()
}
