package utils

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/gob"
	"reflect"
)

func GetBytes(value interface{}) ([]byte, error) {
	if value != nil {
		t := reflect.TypeOf(value)
		v := reflect.New(t).Elem().Interface()
		gob.Register(v)
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GetVal(data []byte, v interface{}) error {
	r := bytes.NewBuffer(data)
	return gob.NewDecoder(r).Decode(v)
}

func ToBase64String(bb []byte) string {
	return b64.StdEncoding.EncodeToString(bb)
}

func FromBase64String(b64str string) ([]byte, error) {
	sDec, err := b64.StdEncoding.DecodeString(b64str)
	if err != nil {
		return nil, err
	}
	return sDec, nil
}
