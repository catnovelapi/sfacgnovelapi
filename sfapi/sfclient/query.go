package sfclient

import "fmt"

type Query map[string]any

func (q Query) Encode() string {
	var s string
	for k, v := range q {
		s += fmt.Sprintf("%v=%v&", k, v)
	}
	return s[:len(s)-1]
}
func (q Query) Add(key string, value any) {
	q[key] = value
}
func (q Query) AddMap(m map[string]any) {
	for k, v := range m {
		q[k] = v
	}

}
func (q Query) ToStringMap() map[string]string {
	m := make(map[string]string)
	for k, v := range q {
		m[k] = fmt.Sprintf("%v", v)
	}
	return m
}
