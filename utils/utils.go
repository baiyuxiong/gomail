package utils

import (
	"strings"
)

func Nl2br(input string) string {
	return strings.Replace(input,"\n","<br/>",-1)
}