package controllers

import (
	"net/http"
	"testing"

	"github.com/gogap/cache_storages"
	"github.com/gogap/spirit"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spirit-contrib/inlet_http"

	"github.com/spirit-contrib/session/defines"
	"github.com/spirit-contrib/session/models"
)

type dataObj struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func TestSessionSet(t *testing.T) {

	var sessionStorage *SessionStorage
	if storage, e := cache_storages.NewMemcachedStorage("127.0.0.1:11211"); e != nil {
		So(e, ShouldBeNil)
	} else {
		sessionStorage = NewSessionStorage(storage)
	}

	Convey("CMD_SESSION_SET command test", t, func() {
		Convey("input correct params, SEESION_KEY not exist before", func() {
			Convey("cache storage will save data correct", func() {
				msg := spirit.Payload{}

				obj := dataObj{Name: "xujinzheng", Email: "xujinzheng@gmail.com"}
				obj2 := dataObj{Name: "xujinzheng2", Email: "xujinzheng2@gmail.com"}

				msg.AppendCommand(defines.CMD_SESSION_SET, models.SessionObject{Name: "sid1", Value: obj})
				msg.AppendCommand(defines.CMD_SESSION_SET, models.SessionObject{Name: "sid2", Value: obj2})

				_, err := sessionStorage.SetSession(&msg)
				So(err, ShouldBeNil)

				cookies := []interface{}{&http.Cookie{}, &http.Cookie{}}
				err = msg.GetCommandObjectArray(inlet_http.CMD_HTTP_COOKIES_SET, cookies)

				So(err, ShouldBeNil)
				So(len(cookies), ShouldEqual, 2)
			})
		})

		Convey("input correct params, SEESION_KEY already exist before", func() {
			Convey("cache storage will get data correct", func() {
				msg := spirit.Payload{}

				obj := dataObj{Name: "xujinzheng", Email: "xujinzheng@gmail.com"}

				ctxCookies := map[string]string{"sid1": "1768f16e-0f33-4e19-58df-0074cbb5376f"}

				msg.SetContext(inlet_http.CTX_HTTP_COOKIES, ctxCookies)

				msg.AppendCommand(defines.CMD_SESSION_SET, models.SessionObject{Name: "sid1", Value: obj})

				_, err := sessionStorage.SetSession(&msg)
				So(err, ShouldBeNil)

				size := msg.GetCommandValueSize(defines.CMD_COOKIES_SET)

				So(err, ShouldBeNil)
				So(size, ShouldEqual, 0)
			})
		})
	})

	Convey("CMD_SESSION_DELETE command test", t, func() {
		Convey("input correct params, SEESION_KEY not exist before", func() {
			Convey("delete SESSION correct", func() {
				sessionStorage.storage.Set("111-111-111", "hello")
				sessionStorage.storage.Set("111-111-222", "world")

				msg := spirit.Payload{}

				cookies := map[string]string{"sid1": "111-111-111", "sid2": "111-111-222"}

				msg.SetContext(inlet_http.CTX_HTTP_COOKIES, cookies)

				msg.AppendCommand(defines.CMD_SESSION_DELETE, "sid1")

				_, err := sessionStorage.SetSession(&msg)
				So(err, ShouldBeNil)

				v1, e1 := sessionStorage.storage.Get("111-111-111")
				So(e1, ShouldNotBeNil)
				So(v1, ShouldEqual, "")

				v2, e2 := sessionStorage.storage.Get("111-111-222")
				So(e2, ShouldBeNil)
				So(v2, ShouldEqual, "world")

				sessionStorage.storage.Delete("111-111-111")
			})
		})
	})
}

func TestSessionQuery(t *testing.T) {
	type dataObj struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	var sessionStorage *SessionStorage
	if storage, e := cache_storages.NewMemcachedStorage("127.0.0.1:11211"); e != nil {
		So(e, ShouldBeNil)
	} else {
		sessionStorage = NewSessionStorage(storage)
	}

	jsonV := `{"v":{"email":"xujinzheng@gmail.com","name":"xujinzheng"}}`
	sessionStorage.storage.Set("1768f16e-0f33-4e19-58df-0074cbb5376f", jsonV)

	Convey("send payload message to query cache", t, func() {
		Convey("input correct params, SEESION_KEY already exist before", func() {
			Convey("cache engine will get session data correct", func() {
				msg := spirit.Payload{}

				ctxCookies := map[string]string{"sid1": "1768f16e-0f33-4e19-58df-0074cbb5376f"}

				msg.SetContext(inlet_http.CTX_HTTP_COOKIES, ctxCookies)

				msg.SetContent(map[string][]string{"cookie_names": []string{"sid1"}})

				result, err := sessionStorage.GetSession(&msg)
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)

				msg2 := spirit.Payload{}
				msg2.SetContent(result)

				v := map[string]dataObj{}
				err = msg2.FillContentToObject(&v)
				So(err, ShouldBeNil)
				So(v["sid1"].Name, ShouldEqual, "xujinzheng")
			})
		})
	})
}
