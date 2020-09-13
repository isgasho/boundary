// Code generated by "make api"; DO NOT EDIT.
package accounts

import "bytes"

type PasswordAccountAttributes struct {
	LoginName string `json:"login_name,omitempty"`
	Password  string `json:"password,omitempty"`

	lastResponseBody *bytes.Buffer
	lastResponseMap  map[string]interface{}
}

func (n PasswordAccountAttributes) LastResponseBody() *bytes.Buffer {
	return n.lastResponseBody
}

func (n PasswordAccountAttributes) LastResponseMap() map[string]interface{} {
	return n.lastResponseMap
}