package keyauth

import (
	"errors"
	"fmt"

	"github.com/rumorshub/gin"
)

var (
	ErrInvalidKey = errors.New("invalid key")
	ErrMissingKey = errors.New("missing key")
)

func Middleware(cfg *Config) gin.HandlerFunc {
	cfg.InitDefaults()

	extractors, err := CreateExtractors(cfg.KeyLookup)
	if err != nil {
		panic(fmt.Errorf("key-auth middleware could not create key extractor: %w", err))
	}
	if len(extractors) == 0 {
		panic(errors.New("key-auth middleware could not create extractors from KeyLookup string"))
	}

	return func(c *gin.Context) {
		rPath := c.Request.Method + c.Request.URL.Path

		for _, rp := range cfg.ExcludeRoutes {
			if gin.MatchPath(rp, rPath) {
				c.Next()
				return
			}
		}

		var (
			lastExtractorErr error
			lastValidatorErr error
		)

		for _, extractor := range extractors {
			keys, source, extrErr := extractor(c)
			if extrErr != nil {
				lastExtractorErr = extrErr
				continue
			}
			for _, key := range keys {
				valid, err := cfg.Validator(c, key, source)
				if err != nil {
					lastValidatorErr = err
					continue
				}
				if !valid {
					lastValidatorErr = gin.ErrUnauthorized().SetInternal(ErrInvalidKey)
					continue
				}

				c.Next()
				return
			}
		}

		err := lastValidatorErr
		if err == nil {
			err = lastExtractorErr
		}

		if cfg.ErrorHandler != nil && err != nil {
			if tmpErr := cfg.ErrorHandler(c, err); tmpErr != nil {
				c.E(tmpErr)
				return
			}

			if cfg.ContinueOnIgnoredError {
				c.Next()
				return
			}
		}

		if lastValidatorErr == nil {
			c.BadRequest(ErrMissingKey)
		} else {
			c.Unauthorized(ErrInvalidKey)
		}
	}
}
