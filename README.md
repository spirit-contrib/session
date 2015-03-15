# session
spirit component - set session, get session, delete session

### Quick start with docker

Configure your ali mqs info in `conf/spirit.env` and set memcached host in `conf/session.conf`

Copy `conf/spirit.env.example` to `conf/spirit.env` and configure it

```json
{
	"owner_id":"",
	"access_key_id": "",
	"acces_key_secert": "",
	"mqs_url":"mqs-cn-hangzhou.aliyuncs.com",
	"queue_get_session":"session-get-session",
	"queue_set_session":"session-set-session"
}
```

Copy `conf/session.conf.example` to `conf/session.conf` and configure it

```json
{
	"memcached":{
		"address":"memcached:11211"
	}
}
```

In these config, we need create two message queue at aliyun console, queue `session-get-session` is for get session from memcached, queue `session-set-session` is for set session to memcached.

> you can also use project of `github.com/spirit-contrib/inlet_http_api` to get session by api method, you just need to do is set `inlet_http_api ` 's graph, and the field of `is_proxy` should be `true`


```bash
$ docker-compose up
Recreating session_memcached_1...
Recreating session_session_1...
Attaching to session_memcached_1, session_session_1
session_1   | [spirit] component session running
```