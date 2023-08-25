package url

import "strings"

var replaceMap = map[string]string{
	"#": "%23",
	"%": "%25",
	"?": "%3F",
	"&": "%26",
}

func Encode(s string) string {
	builder := strings.Builder{}
	builder.Grow(len(s) * 2)
	for _, c := range s {
		if v, ok := replaceMap[string(c)]; ok {
			builder.WriteString(v)
		} else {
			builder.WriteRune(c)
		}
	}
	return builder.String()
}
