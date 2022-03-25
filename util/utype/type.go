package utype

import (
	"fmt"
	"strings"
)

func FullTypeName(value interface{}) string {
	return fmt.Sprintf("%T", value)
}

func TypeName(value interface{}) string {
	parts := strings.Split(FullTypeName(value), ".")
	return parts[len(parts)-1]
}

func AreSameType(v1, v2 interface{}) bool {
	return FullTypeName(v1) == FullTypeName(v2)
}
