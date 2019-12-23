package input

import (
	"bytes"
	"io"
	"log"
	"strconv"
	"strings"
)

var (
	IgnoreCache    bool
	ReportDownload bool
)

func Reader(day int) (io.Reader, error) {
	r, err := daydata(day)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(r), nil
}

func MustReader(day int) io.Reader {
	r, err := Reader(day)
	mustErr(day, err)
	return r
}

func Ints(day int) ([]int, error) {
	data, err := daydata(day)
	if err != nil {
		return nil, err
	}

	s := string(data)
	var sv []string
	if strings.IndexByte(s, ',') > 0 {
		sv = strings.Split(s, ",")
	} else {
		sv = strings.Fields(s)
	}

	var ints []int
	for _, s := range sv {
		n, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			return nil, err
		}
		ints = append(ints, n)
	}

	return ints, nil
}

func MustInts(day int) []int {
	v, err := Ints(day)
	mustErr(day, err)
	return v
}

func String(day int) (string, error) {
	data, err := daydata(day)
	return string(data), err
}

func MustString(day int) string {
	s, err := String(day)
	mustErr(day, err)
	return s
}

func mustErr(day int, err error) {
	if err == nil {
		return
	}

	log.Fatalf("error getting day %d input: %s", day, err)
}
