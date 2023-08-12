package server

import "time"

type Config struct {
	Server      *ConfigServer      `json:"server"`
	Cors        *ConfigCORS        `json:"cors"`
	Compression *ConfigCompression `json:"compression"`
	Limiter     *ConfigLimiter     `json:"limiter"`
	Hosts       []*ConfigHost      `json:"hosts"`
	Static      []*ConfigStatic    `json:"static"`
}

// ConfigServer
// https://dev.to/koddr/go-fiber-by-examples-delving-into-built-in-functions-1p3k
type ConfigServer struct {
	// SETTINGS
	// Request ID adds an identifier to the request using the X-Request-ID header ( uuid.New().String() )
	EnableRequestId bool `json:"enable_request_id"` //
	// Enables use of the SO_REUSEPORT socket option. This will spawn multiple Go processes listening on the same port. learn more about socket sharding
	// https://www.nginx.com/blog/socket-sharding-nginx-release-1-9-1/
	Prefork bool `json:"prefork"` // default: false
	// Enables the Http HTTP header with the given value.
	ServerHeader string `json:"server_header"` // default: ""
	// When enabled, the router treats /foo and /foo/ as different. Otherwise, the router treats /foo and /foo/ as the same.
	StrictRouting bool `json:"strict_routing"` // default: false
	// When enabled, /Foo and /foo are different routes. When disabled, /Fooand /foo are treated the same.
	CaseSensitive bool `json:"case_sensitive"` // default: false
	// When enabled, all values returned by context methods are immutable. By default they are valid until you return from the handler, see issue #185.
	Immutable bool `json:"immutable"` // default: false
	// Sets the maximum allowed size for a request body, if the size exceeds the configured limit, it sends 413 - Request Entity Too Large response.
	BodyLimit int `json:"body_limit"` // default: 4 * 1024 * 1024
	// The amount of time allowed to read the full request including body. Default timeout is unlimited.
	ReadTimeout time.Duration `json:"read_timeout"` // default: 0
	// The maximum duration before timing out writes of the response. Default timeout is unlimited.
	WriteTimeout time.Duration `json:"write_timeout"` // default: 0
	// The maximum amount of time to wait for the next request when keep-alive is enabled. If IdleTimeout is zero, the value of ReadTimeout is used.
	IdleTimeout time.Duration `json:"idle_timeout"` // default: 0
	// Maximum number of concurrent connections.
	Concurrency int `json:"concurrency"` // default: 256 * 1024
	// Disable keep-alive connections. The server will close incoming connections after sending the first response to the client.
	DisableKeepalive bool `json:"disable_keepalive"` // default: false
	// When set to true, it will not print out debug information and startup message
	DisableStartupMessage bool `json:"disable_startup_message"` // default false
}

type ConfigStatic struct {
	Enabled  bool   `json:"enabled"`
	Prefix   string `json:"prefix"`
	Root     string `json:"root"`
	Index    string `json:"index"`
	Compress bool   `json:"compress"`
	// When set to true, enables byte range requests.
	ByteRange bool `json:"byte_range"` // default: false
	// When set to true, enables directory browsing.
	Browse bool `json:"browse"` // default: false.
	// Expiration duration for inactive file handlers. Use a negative time.Duration to disable it
	CacheDurationSec int64 `json:"cache_duration_sec"` // default: 10 * time.Second
	// The value for the Cache-Control HTTP header that is set on the file response. MaxAge is defined in seconds.
	MaxAge int `json:"max_age"` // default: 0

}

type ConfigHost struct {
	Address string `json:"addr"`
	TLS     bool   `json:"tls"`
	// TLS
	SslCert string `json:"ssl_cert"`
	SslKey  string `json:"ssl_key"`
	// websocket
	Websocket *ConfigHostWebsocket `json:"websocket"`
}

type ConfigHostWebsocket struct {
	Enabled bool `json:"enabled"`
	// Specifies the duration for the handshake to complete.
	HandshakeTimeout time.Duration `json:"handshake_timeout"` // default: 0 milliseconds
	// specifies the server's supported protocols in order of preference. If this field is not nil, then the Upgrade
	// method negotiates a subprotocol by selecting the first match in this list with a protocol requested by the client.
	Subprotocols []string `json:"subprotocols"` // default: nil
	// Origins is a string slice of origins that are acceptable, by default all origins are allowed.
	Origins []string `json:"origins"` // default: []string{"*"}
	// ReadBufferSize specify I/O buffer sizes in bytes.
	ReadBufferSize int `json:"read_buffer_size"` // default: 1024
	// WriteBufferSize specify I/O buffer sizes in bytes.
	WriteBufferSize int `json:"write_buffer_size"` // default: 1024
	// EnableCompression specify if the server should attempt to negotiate per message compression (RFC 7692)
	EnableCompression bool `json:"enable_compression"` // default:false
}

type ConfigCORS struct {
	Enabled bool `json:"enabled"`
	// AllowOrigin defines a list of origins that may access the resource.
	AllowOrigins []string `json:"allow_origins"` // default: []string{"*"}
	// AllowMethods defines a list methods allowed when accessing the resource. This is used in response to a preflight request.
	AllowMethods []string `json:"allow_methods"` // default: []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"}
	// AllowCredentials indicates whether or not the response to the request can be exposed when the credentials flag is true. When used as part of a response to a preflight request, this indicates whether or not the actual request can be made using credentials.
	AllowCredentials bool `json:"allow_credentials"` // default: false
	// ExposeHeaders defines a whitelist headers that clients are allowed to access.
	ExposeHeaders []string `json:"expose_headers"` // default: nil
	// MaxAge indicates how long (in seconds) the results of a preflight request can be cached.
	MaxAge int `json:"max_age"` // default: 0
}

type ConfigCompression struct {
	Enabled bool `json:"enabled"`
	Level   int  `json:"level"` // Level of compression, 0, 1, 2, 3, 4
}

type ConfigLimiter struct {
	Enabled bool `json:"enabled"`
	// Max number of recent connections during `Duration` seconds before sending a 429 response
	//
	// Default: 5
	Max int `json:"max"`

	// Duration is the time on how long to keep records of requests in memory
	//
	// Default: 1 * time.Minute
	Duration time.Duration `json:"duration"`
}
