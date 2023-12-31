package json

import (
	"encoding/json"
	"io"
	"os"

	jsoniter "github.com/json-iterator/go"
)

var j = jsoniter.ConfigCompatibleWithStandardLibrary

type Number = json.Number

type RawMessage = json.RawMessage

func Marshal(v any) ([]byte, error) {
	return j.Marshal(v)
}

func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return j.MarshalIndent(v, prefix, indent)
}

func MarshalString(v any) string {
	by, err := j.Marshal(v)
	if err != nil {
		return ""
	}
	return string(by)
}

func MarshalStringPretty(v any) string {
	by, err := MarshallPretty(v)
	if err != nil {
		return ""
	}
	return string(by)
}

func MarshallPretty(v any) ([]byte, error) {
	return j.MarshalIndent(v, "", "    ")
}

func Unmarshal(data []byte, v any) error {
	return j.Unmarshal(data, v)
}

func UnmarshalString(str string, v interface{}) error {
	return j.UnmarshalFromString(str, v)
}

func UnmarshalFile(path string, v interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	by, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	return j.Unmarshal(by, v)
}

func Get(data []byte, path ...interface{}) jsoniter.Any {
	return j.Get(data, path...)
}

func Valid(data []byte) bool {
	return j.Valid(data)
}

func NewDecoder(r io.Reader) *jsoniter.Decoder {
	return j.NewDecoder(r)
}

func NewEncoder(w io.Writer) *jsoniter.Encoder {
	return j.NewEncoder(w)
}
