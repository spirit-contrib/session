package controllers

import (
	"net/http"
	"strings"

	"github.com/gogap/cache_storages"
	"github.com/gogap/errors"
	"github.com/gogap/logs"
	"github.com/gogap/spirit"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/spirit-contrib/inlet_http"

	"github.com/spirit-contrib/session/defines"
	"github.com/spirit-contrib/session/errorcode"
	"github.com/spirit-contrib/session/models"
)

type SessionStorage struct {
	serverSalt string
	storage    cache_storages.CacheStorage
}

func NewSessionStorage(storage cache_storages.CacheStorage) *SessionStorage {
	newStorage := new(SessionStorage)
	newStorage.storage = storage
	return newStorage
}

func (p *SessionStorage) SetSession(msg *spirit.Payload) (result interface{}, err error) {
	result = msg.GetContent()

	if size := msg.GetCommandValueSize(defines.CMD_SESSION_SET); size > 0 {
		var values []interface{}
		values = make([]interface{}, size)
		for i := 0; i < size; i++ {
			values[i] = &models.SessionObject{}
		}

		//GetCommand from pre component
		if e := msg.GetCommandObjectArray(defines.CMD_SESSION_SET, values); e != nil {
			err = errorcode.ERR_SESSION_PARSE_CMD_TO_OBJ_FAILED.New(
				errors.Params{"err": e})
			logs.Error(err)
			return
		}

		cookies := map[string]string{}
		msg.GetContextObject(inlet_http.CTX_HTTP_COOKIES, &cookies)

		//SetSessionData to Memcached
		for _, ikvs := range values {
			sObj, _ := ikvs.(*models.SessionObject)

			if strings.TrimSpace(sObj.Name) == "" {
				err = errorcode.ERR_SESSION_NAME_IS_EMPTY.New()
				return
			}

			//Check if user already have cookies
			userCookieId := ""

			if cookieId, exist := cookies[sObj.Name]; !exist {
				//Set Command to next component
				if id, e := uuid.NewV4(); e != nil {
					err = errorcode.ERR_SESSION_GENERATE_ID_FAILED.New(errors.Params{"err": e})
					return
				} else {
					userCookieId = id.String()
				}
				newCookie := http.Cookie{
					Name:     sObj.Name,
					Value:    userCookieId,
					HttpOnly: true,
				}

				logs.Debug("set new cookie for seesion id:", userCookieId)
				msg.AppendCommand(inlet_http.CMD_HTTP_COOKIES_SET, newCookie)
			} else {
				userCookieId = cookieId
			}
			p.storage.SetObject(userCookieId, sObj.Value)
		}
	}

	if values, e := msg.GetCommandStringArray(defines.CMD_SESSION_DELETE); e == nil && values != nil && len(values) > 0 {
		cookies := map[string]string{}
		msg.GetContextObject(inlet_http.CTX_HTTP_COOKIES, &cookies)

		//SetSessionData to Memcached
		for _, cookieName := range values {
			if cookieId, exist := cookies[cookieName]; exist {
				p.storage.Delete(cookieId)
			} else {
				logs.Debug("the session to be delete did not have related cookie id:", cookieName)
			}
		}
	}
	return
}

func (p *SessionStorage) GetSession(msg *spirit.Payload) (result interface{}, err error) {
	queryCookieNames := []string{}

	values := map[string]interface{}{}

	if e := msg.FillContentToObject(&queryCookieNames); e != nil {
		err = errorcode.ERR_COOKIES_GET_NAME_FAILED.New(errors.Params{"err": e})
		return
	} else {
		cookies := map[string]string{}
		msg.GetContextObject(inlet_http.CTX_HTTP_COOKIES, &cookies)

		keyMap := map[string]string{}
		for _, queryName := range queryCookieNames {
			if sid, exist := cookies[queryName]; exist {
				values[sid] = &map[string]interface{}{}
				keyMap[sid] = queryName
			}
		}

		if e := p.storage.GetMultiObject(values); e != nil {
			err = errorcode.ERR_SESSION_GET_DATA_FAILED.New(errors.Params{"err": e})
			logs.Error(err)
			return
		}

		ret := map[string]interface{}{}
		for sid, value := range values {
			name, _ := keyMap[sid]
			ret[name] = value
		}
		result = ret
	}
	return
}
