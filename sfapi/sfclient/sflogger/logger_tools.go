package sflogger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
)

func indentJson(a string) string {
	var objmap map[string]*json.RawMessage
	err := json.Unmarshal([]byte(a), &objmap)
	if err != nil {
		log.Println(err)
		return a + "\n" + err.Error()
	}
	formatted, err := json.MarshalIndent(objmap, "", "  ")
	if err != nil {
		return a + "\n" + err.Error()
	}
	return string(formatted)
}

func copyHeaders(hdrs http.Header) http.Header {
	nh := http.Header{}
	if hdrs != nil {
		for k, v := range hdrs {
			nh[k] = v
		}
	}
	return nh
}
func sortHeaderKeys(hdrs http.Header) []string {
	keys := make([]string, 0, len(hdrs))
	for key := range hdrs {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func composeHeaders(hdrs http.Header) string {
	str := make([]string, 0, len(hdrs))
	for _, k := range sortHeaderKeys(hdrs) {
		str = append(str, "\t"+strings.TrimSpace(fmt.Sprintf("%25s: %s", k, strings.Join(hdrs[k], ", "))))
	}
	return strings.Join(str, "\n")
}
