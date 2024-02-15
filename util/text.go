package util

import (
	"fmt"
	"strings"
)

func Initcap(s string) string {
	return fmt.Sprint(strings.ToUpper(s[:1]), s[1:])
}
