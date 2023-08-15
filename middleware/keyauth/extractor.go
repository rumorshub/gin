package keyauth

import (
	"fmt"
	"net/textproto"
	"strings"

	"github.com/rumorshub/gin"
)

const extractorLimit = 20

type ExtractorSource string

const (
	ExtractorSourceHeader ExtractorSource = "header"
	ExtractorSourceQuery  ExtractorSource = "query"
	ExtractorSourcePath   ExtractorSource = "path"
	ExtractorSourceForm   ExtractorSource = "form"
	ExtractorSourceCtx    ExtractorSource = "ctx"
)

type ValueExtractorError struct {
	message string
}

func (e *ValueExtractorError) Error() string {
	return e.message
}

var (
	ErrHeaderExtractorValueMissing = &ValueExtractorError{message: "missing value in request header"}
	ErrHeaderExtractorValueInvalid = &ValueExtractorError{message: "invalid value in request header"}
	ErrQueryExtractorValueMissing  = &ValueExtractorError{message: "missing value in query string"}
	ErrPathExtractorValueMissing   = &ValueExtractorError{message: "missing value in path params"}
	ErrFormExtractorValueMissing   = &ValueExtractorError{message: "missing value in form"}
	ErrCtxExtractorValueMissing    = &ValueExtractorError{message: "missing value in ctx"}
)

type ValuesExtractor func(c *gin.Context) ([]string, ExtractorSource, error)

func CreateExtractors(lookups string) ([]ValuesExtractor, error) {
	if lookups == "" {
		return nil, nil
	}
	sources := strings.Split(lookups, ",")
	var extractors = make([]ValuesExtractor, 0)
	for _, source := range sources {
		parts := strings.Split(source, ":")
		if len(parts) < 2 {
			return nil, fmt.Errorf("extractor source for lookup could not be split into needed parts: %v", source)
		}

		switch parts[0] {
		case "query":
			extractors = append(extractors, ValuesFromQuery(parts[1]))
		case "path":
			extractors = append(extractors, ValuesFromPath(parts[1]))
		case "ctx":
			extractors = append(extractors, ValuesFromCtx(parts[1]))
		case "form":
			extractors = append(extractors, ValuesFromForm(parts[1]))
		case "header":
			prefix := ""
			if len(parts) > 2 {
				prefix = parts[2]
			}
			extractors = append(extractors, ValuesFromHeader(parts[1], prefix))
		}
	}
	return extractors, nil
}

func ValuesFromHeader(header string, valuePrefix string) ValuesExtractor {
	prefixLen := len(valuePrefix)
	header = textproto.CanonicalMIMEHeaderKey(header)
	return func(c *gin.Context) ([]string, ExtractorSource, error) {
		values := c.Request.Header.Values(header)
		if len(values) == 0 {
			return nil, ExtractorSourceHeader, ErrHeaderExtractorValueMissing
		}

		var result []string
		for i, value := range values {
			if prefixLen == 0 {
				result = append(result, value)
				if i >= extractorLimit-1 {
					break
				}
				continue
			}
			if len(value) > prefixLen && strings.EqualFold(value[:prefixLen], valuePrefix) {
				result = append(result, value[prefixLen:])
				if i >= extractorLimit-1 {
					break
				}
			}
		}

		if len(result) == 0 {
			if prefixLen > 0 {
				return nil, ExtractorSourceHeader, ErrHeaderExtractorValueInvalid
			}
			return nil, ExtractorSourceHeader, ErrHeaderExtractorValueMissing
		}
		return result, ExtractorSourceHeader, nil
	}
}

func ValuesFromQuery(param string) ValuesExtractor {
	return func(c *gin.Context) ([]string, ExtractorSource, error) {
		if result, ok := c.GetQueryArray(param); ok && len(result) > 0 {
			return result, ExtractorSourceQuery, nil
		}
		return nil, ExtractorSourceQuery, ErrQueryExtractorValueMissing
	}
}

func ValuesFromPath(param string) ValuesExtractor {
	return func(c *gin.Context) ([]string, ExtractorSource, error) {
		if result := c.Param(param); result != "" {
			return []string{result}, ExtractorSourcePath, nil
		}
		return nil, ExtractorSourcePath, ErrPathExtractorValueMissing
	}
}

func ValuesFromCtx(name string) ValuesExtractor {
	return func(c *gin.Context) ([]string, ExtractorSource, error) {
		if result := c.GetString(name); result != "" {
			return []string{result}, ExtractorSourceCtx, nil
		}
		return nil, ExtractorSourceCtx, ErrCtxExtractorValueMissing
	}
}

func ValuesFromForm(name string) ValuesExtractor {
	return func(c *gin.Context) ([]string, ExtractorSource, error) {
		if result, ok := c.GetPostFormArray(name); ok && len(result) > 0 {
			return result, ExtractorSourceForm, nil
		}
		return nil, ExtractorSourceForm, ErrFormExtractorValueMissing
	}
}
