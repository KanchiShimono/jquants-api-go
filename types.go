package jquants_api_go

import (
	"bytes"
	"fmt"
	"time"
)

const (
	codeQueryKey = "code"
	dateQueryKey = "date"
)

type JSONTime int64

/*
https://kenzo0107.github.io/2020/05/19/2020-05-20-go-json-time/
*/

// String converts the unix timestamp into a string
func (t JSONTime) String() string {
	tm := t.Time()
	return fmt.Sprintf("\"%s\"", tm.Format("2006-01-02"))
}

// Time returns a `time.Time` representation of this value.
func (t JSONTime) Time() time.Time {
	return time.Unix(int64(t), 0)
}

// UnmarshalJSON will unmarshal both string and int JSON values
func (t *JSONTime) UnmarshalJSON(buf []byte) error {
	s := bytes.Trim(buf, `"`)
	aa, err := time.Parse("20060102", string(s))
	if err != nil {
		return err
	}

	*t = JSONTime(aa.Unix())
	return nil
}
