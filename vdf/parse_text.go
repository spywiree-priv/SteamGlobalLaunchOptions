package vdf

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"unsafe"

	"github.com/guregu/null/v5"
)

var re = regexp.MustCompile(`(\".*?\")(?:\t\t(\".*\"))?`)

func bytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func ParseText(r io.Reader) (*KeyValue, error) {
	var root KeyValue
	level := 0

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := bytes.TrimSpace(sc.Bytes())
		if len(line) == 0 || bytes.HasPrefix(line, []byte("/")) {
			continue
		}

		if bytes.Equal(line, []byte("{")) {
			level++
			continue
		} else if bytes.Equal(line, []byte("}")) {
			level--
			continue
		}

		match := re.FindStringSubmatch(bytesToString(line))
		if len(match) < 3 {
			continue
		}
		key := Unquote(match[1])
		var value null.String
		if match[2] != "" {
			value.SetValid(Unquote(match[2]))
		}

		if level == 0 {
			root.Key = key
			root.Value = value
			continue
		}

		parent := &root
		for range level - 1 {
			if len(parent.Children) == 0 {
				parent.Children = append(parent.Children, KeyValue{})
			}
			parent = &parent.Children[len(parent.Children)-1]
		}
		parent.Children = append(parent.Children, KeyValue{Key: key, Value: value})
	}

	return &root, nil
}
