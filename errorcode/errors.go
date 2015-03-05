package errorcode

import (
	"github.com/gogap/errors"
)

const (
	SESSION_ERROR_NS = "SESSION"
)

var (
	ERR_SESSION_PARSE_CMD_TO_OBJ_FAILED = errors.TN(SESSION_ERROR_NS, 1, "parse session cmd values to object failed, raw error is {{.err}}")
	ERR_SESSION_SET_DATA_FAILED         = errors.TN(SESSION_ERROR_NS, 2, "set session data failed, key: {{.key}}, value: {{.value}}, raw error is {{.err}}")
	ERR_SESSION_GENERATE_ID_FAILED      = errors.TN(SESSION_ERROR_NS, 3, "generate session id failed")
	ERR_SESSION_GET_COOKIE_ID_FAILED    = errors.TN(SESSION_ERROR_NS, 4, "get get cookies id failed, raw error is: {{.err}}")
	ERR_SESSION_NAME_IS_EMPTY           = errors.TN(SESSION_ERROR_NS, 5, "session name is empty")
	ERR_SESSION_GET_DATA_FAILED         = errors.TN(SESSION_ERROR_NS, 6, "get session data failed")
	ERR_COOKIES_GET_NAME_FAILED         = errors.TN(SESSION_ERROR_NS, 7, "get cookies name list faild")
)
