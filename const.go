package skynology

// Skynology api methods.
const (
	GET    Method = "GET"
	POST   Method = "POST"
	DELETE Method = "DELETE"
	PUT    Method = "PUT"
)

const (
	X_APPLICATION_ID_HEADER = "X-Sky-Application-Id"
	X_REQUEST_SIGN_HEADER   = "X-Sky-Request-Sign"
	X_SESSION_TOKEN_HEADER  = "X-Sky-Session-Token"
	X_CLIENT_VERSION_HEADER = "X-Sky-Client-Version"
	X_WEIXIN_ID_HEADER      = "X-Sky-Weixin-Id"
	X_WEIXIN_TYPE_HEADER    = "X-Sky-Weixin-Type"
)

const (
	SDK_VERSION = "0.1.0"
)
