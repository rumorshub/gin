package gin

type Config struct {
	// Mode gin mode according to input string: debug, release, test.
	Mode string `mapstructure:"mode" json:"mode,omitempty" bson:"mode,omitempty"`

	// Static files which serves from the given file system root.
	Static map[string]string `mapstructure:"static" json:"static,omitempty" bson:"static,omitempty"`

	// List of the middleware names (order will be preserved).
	Middleware []string `mapstructure:"middleware" json:"middleware,omitempty" bson:"middleware,omitempty"`

	// EnableDecoderUseNumber is used to call the UseNumber method on the JSON
	// Decoder instance. UseNumber causes the Decoder to unmarshal a number into an
	// any as a Number instead of as a float64.
	EnableDecoderUseNumber bool `mapstructure:"enable_decoder_use_number" json:"enable_decoder_use_number,omitempty" bson:"enable_decoder_use_number,omitempty"`

	// EnableDecoderDisallowUnknownFields is used to call the DisallowUnknownFields method
	// on the JSON Decoder instance. DisallowUnknownFields causes the Decoder to
	// return an error when the destination is a struct and the input contains object
	// keys which do not match any non-ignored, exported fields in the destination.
	EnableDecoderDisallowUnknownFields bool `mapstructure:"enable_decoder_disallow_unknown_fields" json:"enable_decoder_disallow_unknown_fields,omitempty" bson:"enable_decoder_disallow_unknown_fields,omitempty"`

	// RedirectTrailingSlash enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 307 for all other request methods.
	RedirectTrailingSlash *bool `mapstructure:"redirect_trailing_slash" json:"redirect_trailing_slash,omitempty" bson:"redirect_trailing_slash,omitempty"`

	// RedirectFixedPath if enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 307 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool `mapstructure:"redirect_fixed_path" json:"redirect_fixed_path,omitempty" bson:"redirect_fixed_path,omitempty"`

	// HandleMethodNotAllowed if enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool `mapstructure:"handle_method_not_allowed" json:"handle_method_not_allowed,omitempty" bson:"handle_method_not_allowed,omitempty"`

	// ForwardedByClientIP if enabled, client IP will be parsed from the request's headers that
	// match those stored at `(*gin.Engine).RemoteIPHeaders`. If no IP was
	// fetched, it falls back to the IP obtained from
	// `(*gin.Context).Request.RemoteAddr`
	ForwardedByClientIP *bool `mapstructure:"forwarded_by_client_ip" json:"forwarded_by_client_ip,omitempty" bson:"forwarded_by_client_ip,omitempty"`

	// UseRawPath if enabled, the url.RawPath will be used to find parameters.
	UseRawPath bool `mapstructure:"use_raw_path" json:"use_raw_path,omitempty" bson:"use_raw_path,omitempty"`

	// UnescapePathValues if true, the path value will be unescaped.
	// If UseRawPath is false (by default), the UnescapePathValues effectively is true,
	// as url.Path gonna be used, which is already unescaped.
	UnescapePathValues *bool `mapstructure:"unescape_path_values" json:"unescape_path_values,omitempty" bson:"unescape_path_values,omitempty"`

	// RemoveExtraSlash a parameter can be parsed from the URL even with extra slashes.
	// See the PR #1817 and issue #1644
	RemoveExtraSlash bool `mapstructure:"remove_extra_slash" json:"remove_extra_slash,omitempty" bson:"remove_extra_slash,omitempty"`

	// RemoteIPHeaders list of headers used to obtain the client IP when
	// `(*gin.Engine).ForwardedByClientIP` is `true` and
	// `(*gin.Context).Request.RemoteAddr` is matched by at least one of the
	// network origins of list defined by `(*gin.Engine).SetTrustedProxies()`.
	RemoteIPHeaders []string `mapstructure:"remote_ip_headers" json:"remote_ip_headers,omitempty" bson:"remote_ip_headers,omitempty"`

	// TrustedPlatform if set to a constant of value gin.Platform*, trusts the headers set by
	// that platform, for example to determine the client IP
	TrustedPlatform string `mapstructure:"trusted_platform" json:"trusted_platform,omitempty" bson:"trusted_platform,omitempty"`

	// MaxMultipartMemory value of 'maxMemory' param that is given to http.Request's ParseMultipartForm
	// method call.
	MaxMultipartMemory int64 `mapstructure:"max_multipart_memory" json:"max_multipart_memory,omitempty" bson:"max_multipart_memory,omitempty"`

	// ContextWithFallback enable fallback Context.Deadline(), Context.Done(), Context.Err() and Context.Value() when Context.Request.Context() is not nil.
	ContextWithFallback bool `mapstructure:"context_with_fallback" json:"context_with_fallback,omitempty" bson:"context_with_fallback,omitempty"`
}

func (g *Config) InitDefaults() {
	if g.Mode == "" {
		g.Mode = "release"
	}
	if g.RedirectTrailingSlash == nil {
		g.RedirectTrailingSlash = toPtr(true)
	}
	if g.ForwardedByClientIP == nil {
		g.ForwardedByClientIP = toPtr(true)
	}
	if g.UnescapePathValues == nil {
		g.UnescapePathValues = toPtr(true)
	}
	if len(g.RemoteIPHeaders) == 0 {
		g.RemoteIPHeaders = []string{"X-Forwarded-For", "X-Real-IP"}
	}
	if g.MaxMultipartMemory == 0 {
		g.MaxMultipartMemory = 32 << 20 // 32 MB
	}
}

func toPtr[T any](val T) *T {
	return &val
}
