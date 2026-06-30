package tmpl

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"strings"
)

// FuncMap is the standard library of template helpers available to every
// embedded template. It grows as templates need more; everything here is
// deterministic so incremental builds stay correct.
func FuncMap() template.FuncMap {
	return template.FuncMap{
		// Defaults / selection.
		"default": func(def any, val any) any {
			if isEmpty(val) {
				return def
			}
			return val
		},
		"coalesce": func(vals ...any) any {
			for _, v := range vals {
				if !isEmpty(v) {
					return v
				}
			}
			return ""
		},

		// Numeric conversion + arithmetic (args arrive as strings).
		"int":   toInt,
		"float": toFloat,
		"add":   func(a, b any) float64 { return toFloat(a) + toFloat(b) },
		"sub":   func(a, b any) float64 { return toFloat(a) - toFloat(b) },
		"mul":   func(a, b any) float64 { return toFloat(a) * toFloat(b) },
		"div": func(a, b any) float64 {
			if toFloat(b) == 0 {
				return 0
			}
			return toFloat(a) / toFloat(b)
		},
		"seq": seq,

		// Strings.
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": strings.Title,
		"trim":  strings.TrimSpace,
		"split": func(sep, s string) []string { return strings.Split(s, sep) },
		"join":  func(sep string, parts []string) string { return strings.Join(parts, sep) },
		"replace": func(old, new, s string) string {
			return strings.ReplaceAll(s, old, new)
		},
		"contains": strings.Contains,

		// Data construction, useful for passing structured config into templates.
		"dict": dict,
		"list": func(vals ...any) []any { return vals },

		// JSON encoding for embedding config in <script> blocks.
		"json": func(v any) (template.JS, error) {
			b, err := json.Marshal(v)
			return template.JS(b), err
		},
		"jsonIndent": func(v any) (template.JS, error) {
			b, err := json.MarshalIndent(v, "", "  ")
			return template.JS(b), err
		},

		// base64 (URL-safe, unpadded) — carries JavaScript payloads (e.g. sketch
		// source) through HTML untouched. The URL-safe alphabet avoids '+' and
		// '/', which html/template escapes inside a <script> element.
		"b64": func(s string) string { return base64.RawURLEncoding.EncodeToString([]byte(s)) },

		// Escape hatches: opt into trusting a value in a specific context.
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
		"safeJS":   func(s string) template.JS { return template.JS(s) },
		"safeCSS":  func(s string) template.CSS { return template.CSS(s) },
		"safeURL":  func(s string) template.URL { return template.URL(s) },
		"attr":     func(s string) template.HTMLAttr { return template.HTMLAttr(s) },

		"bool": func(v any) bool { return truthy(v) },
	}
}

func dict(pairs ...any) (map[string]any, error) {
	if len(pairs)%2 != 0 {
		return nil, fmt.Errorf("dict: odd number of arguments (%d)", len(pairs))
	}
	m := make(map[string]any, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		key, ok := pairs[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict: key %d is not a string", i)
		}
		m[key] = pairs[i+1]
	}
	return m, nil
}

func seq(args ...any) []int {
	start, end, step := 0, 0, 1
	switch len(args) {
	case 1:
		end = toInt(args[0])
	case 2:
		start, end = toInt(args[0]), toInt(args[1])
	case 3:
		start, end, step = toInt(args[0]), toInt(args[1]), toInt(args[2])
	default:
		return nil
	}
	if step == 0 {
		return nil
	}
	var out []int
	if step > 0 {
		for i := start; i < end; i += step {
			out = append(out, i)
		}
	} else {
		for i := start; i > end; i += step {
			out = append(out, i)
		}
	}
	return out
}

func toInt(v any) int {
	switch x := v.(type) {
	case int:
		return x
	case float64:
		return int(x)
	case string:
		n, _ := strconv.Atoi(strings.TrimSpace(x))
		return n
	}
	return 0
}

func toFloat(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	case string:
		f, _ := strconv.ParseFloat(strings.TrimSpace(x), 64)
		return f
	}
	return 0
}

func truthy(v any) bool {
	switch x := v.(type) {
	case bool:
		return x
	case string:
		s := strings.ToLower(strings.TrimSpace(x))
		return s == "true" || s == "1" || s == "yes" || s == "on"
	case int:
		return x != 0
	case float64:
		return x != 0
	}
	return v != nil
}

func isEmpty(v any) bool {
	switch x := v.(type) {
	case nil:
		return true
	case string:
		return x == ""
	case int:
		return x == 0
	case float64:
		return x == 0
	case bool:
		return !x
	case []any:
		return len(x) == 0
	}
	return false
}
