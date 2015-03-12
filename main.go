package main

import (
	"encoding/json"
	"errors"
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

	sessionStorage := new(controllers.SessionStorage)

	funcInitalSession := func(configFile string) (err error) {
		logs.SetFileLogger("./logs/session.log")

		if configFile == "" {
			configFile = "./conf/session.conf"
		}

		storageAddr := ""
		if bFile, e := ioutil.ReadFile(configFile); e != nil {
			err = e
			return
		} else {
			sessionConfig := SessionConfig{}
			if e := json.Unmarshal(bFile, &sessionConfig); e != nil {
				err = e
				return
			}

			sessionConfig.Address = strings.TrimSpace(sessionConfig.Address)
			if sessionConfig.Address == "" {
				err = errors.New("memcached.address is empty")
				return
			}

			storageAddr = sessionConfig.Address
		}

		if storage, e := cache_storages.NewMemcachedStorage(storageAddr); e != nil {
			err = e
			return
		} else {
			sessionStorage.SetStorage(storage)
		}

		return
	}

	sessionSpirit := spirit.NewClassicSpirit("session", "a basic session component", "1.0.0")
	sessionComponent := spirit.NewBaseComponent("session")

	sessionComponent.RegisterHandler("set_session", sessionStorage.SetSession)
	sessionComponent.RegisterHandler("get_session", sessionStorage.GetSession)

	sessionSpirit.Hosting(sessionComponent).Build().Run(funcInitalSession)
}
