package main

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/gogap/cache_storages"
	"github.com/gogap/env_json"
	"github.com/gogap/logs"
	"github.com/gogap/spirit"

	"github.com/spirit-contrib/session/controllers"
)

const (
	SESSION_ENV = "SESSION_ENV"
)

type SessionConfig struct {
	MemcachedConfig MemcachedConfig `json:"memcached"`
	SessionConf     SessionConf     `json:"session"`
}

type MemcachedConfig struct {
	Address string `json:"address"`
}

type SessionConf struct {
	ExpirationSeconds int32 `json:"expiration_seconds"`
}

func main() {

	sessionStorage := new(controllers.SessionStorage)

	funcInitalSession := func() (err error) {
		logs.SetFileLogger("logs/session.log")

		storageAddr := ""
		if bFile, e := ioutil.ReadFile("conf/session.conf"); e != nil {
			err = e
			return
		} else {
			sessionConfig := SessionConfig{}
			envJson := env_json.NewEnvJson(SESSION_ENV, env_json.ENV_JSON_EXT)
			if e := envJson.Unmarshal(bFile, &sessionConfig); e != nil {
				err = e
				return
			}

			sessionConfig.MemcachedConfig.Address = strings.TrimSpace(sessionConfig.MemcachedConfig.Address)
			if sessionConfig.MemcachedConfig.Address == "" {
				err = errors.New("memcached.address is empty")
				return
			}

			if sessionConfig.SessionConf.ExpirationSeconds < 0 {
				sessionConfig.SessionConf.ExpirationSeconds = 0
			}

			sessionStorage.SetExpireSeconds(sessionConfig.SessionConf.ExpirationSeconds)

			logs.Info("session expire seconds:", sessionConfig.SessionConf.ExpirationSeconds)

			storageAddr = sessionConfig.MemcachedConfig.Address
		}

		if storage, e := cache_storages.NewMemcachedStorage(storageAddr); e != nil {
			err = e
			return
		} else {
			sessionStorage.SetStorage(storage)
		}

		return
	}

	sessionSpirit := spirit.NewClassicSpirit(
		"session",
		"a basic session component",
		"1.0.0",
		[]spirit.Author{
			{Name: "zeal", Email: "xujinzheng@gmail.com"},
		})

	sessionComponent := spirit.NewBaseComponent("session")

	sessionComponent.RegisterHandler("set_session", sessionStorage.SetSession)
	sessionComponent.RegisterHandler("get_session", sessionStorage.GetSession)

	sessionSpirit.Hosting(sessionComponent, funcInitalSession).Build().Run()
}
