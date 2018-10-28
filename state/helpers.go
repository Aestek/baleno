package state

import (
	"fmt"
	"strings"
)

func K(parts ...interface{}) string {
	key := strings.Builder{}
	for i, p := range parts {
		if i > 0 {
			key.WriteByte('.')
		}
		key.WriteString(fmt.Sprint(p))
	}
	return key.String()
}
