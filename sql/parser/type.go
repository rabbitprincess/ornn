package parser

import (
	"regexp"
	"strconv"
	"strings"
)

// Type holds information for a database type.
type Type struct {
	Type     string `json:"type,omitempty"`
	Prec     int    `json:"prec,omitempty"`
	Scale    int    `json:"scale,omitempty"`
	Nullable bool   `json:"nullable,omitempty"`
	IsArray  bool   `json:"array,omitempty"`
	Unsigned bool   `json:"unsigned,omitempty"`
}

func ParseType(typ string) *Type {
	isArray := false
	if strings.HasSuffix(typ, "[]") {
		typ, isArray = typ[:len(typ)-len("[]")], true
	}
	unsigned := false
	if strings.HasSuffix(typ, " unsigned") {
		typ, unsigned = typ[:len(typ)-len(" unsigned")], true
	}
	var prec, scale int
	if m := regexp.MustCompile(`\(([0-9]+)(\s*,\s*[0-9]+\s*)?\)$`).FindStringIndex(typ); m != nil {
		s := typ[m[0]+1 : m[1]-1]
		if i := strings.LastIndex(s, ","); i != -1 {
			scale, _ = strconv.Atoi(strings.TrimSpace(s[i+1:]))
			s = s[:i]
		}
		prec, _ = strconv.Atoi(strings.TrimSpace(s))
		typ = typ[:m[0]]
	}
	return &Type{
		Type:     strings.ToLower(strings.TrimSpace(typ)),
		Prec:     prec,
		Scale:    scale,
		IsArray:  isArray,
		Unsigned: unsigned,
	}
}
