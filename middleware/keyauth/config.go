package keyauth

import (
	"errors"

	"github.com/rumorshub/gin"
)

// Validator defines a function to validate KeyAuth credentials.
type Validator func(c *gin.Context, key string, source ExtractorSource) (bool, error)

// ErrorHandler defines a function which is executed for an invalid key.
type ErrorHandler func(c *gin.Context, err error) error

type Config struct {
	// KeyLookup is a string in the form of "<source>:<name>" or "<source>:<name>,<source>:<name>" that is used
	// to extract key from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>" or "header:<name>:<cut-prefix>"
	// 			`<cut-prefix>` is argument value to cut/trim prefix of the extracted value. This is useful if header
	//			value has static prefix like `Authorization: <auth-scheme> <authorisation-parameters>` where part that we
	//			want to cut is `<auth-scheme> ` note the space at the end.
	//			In case of basic authentication `Authorization: Basic <credentials>` prefix we want to remove is `Basic `.
	// - "query:<name>"
	// - "form:<name>"
	// - "cookie:<name>"
	// Multiple sources example:
	// - "header:Authorization,header:X-Api-Key"
	KeyLookup string `mapstructure:"key_lookup"`

	// ContinueOnIgnoredError allows the next middleware/handler to be called when ErrorHandler decides to
	// ignore the error (by returning `nil`).
	// This is useful when parts of your site/api allow public access and some authorized routes provide extra functionality.
	// In that case you can use ErrorHandler to set a default public key auth value in the request context
	// and continue. Some logic down the remaining execution chain needs to check that (public) key auth value then.
	ContinueOnIgnoredError bool `mapstructure:"continue_on_ignored_error"`

	ExcludeRoutes []string `mapstructure:"exclude_routes"`

	// Validator is a function to validate key.
	// Required.
	Validator Validator

	// ErrorHandler defines a function which is executed for an invalid key.
	// It may be used to define a custom error.
	ErrorHandler ErrorHandler
}

func (cfg *Config) InitDefaults() {
	if cfg.Validator == nil {
		panic(errors.New("key-auth middleware requires a validator function"))
	}
	if cfg.KeyLookup == "" {
		cfg.KeyLookup = "header:Authorization:Bearer "
	}
}
