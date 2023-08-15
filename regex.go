package gin

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

var regPath = new(sync.Map)

func GetRePath(str string) *regexp.Regexp {
	if v, ok := regPath.Load(str); ok {
		return v.(*regexp.Regexp)
	}

	tmp := strings.ReplaceAll(str, "*", "(|.+)")
	re := regexp.MustCompile(fmt.Sprintf("(?i)^%s$", tmp))
	regPath.Store(str, re)

	return re
}

func MatchPath(regPath, rPath string) bool {
	return GetRePath(regPath).MatchString(rPath)
}
