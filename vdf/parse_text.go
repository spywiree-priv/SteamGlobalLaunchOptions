package vdf

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`\"(.*?)\"(?:\t\t\"(.*)\")?`)

func ParseText(r io.Reader) (KeyValue, error) {
	var root KeyValue
	level := 0

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "/") {
			continue
		}

		if line == "{" {
			level += 1
			continue
		} else if line == "}" {
			level -= 1
			continue
		}

		match := re.FindStringSubmatch(line)
		var key, value string
		if len(match) >= 2 {
			key = match[1]
		}
		if len(match) >= 3 {
			value = match[2]
		}

		if level == 0 {
			root.Key = key
			root.Value = value
			continue
		}

		parent := &root
		for i := 0; i < level-1; i++ {
			if len(parent.Children) == 0 {
				parent.Children = append(parent.Children, KeyValue{})
			}
			parent = &parent.Children[len(parent.Children)-1]
		}
		parent.Children = append(parent.Children, KeyValue{Key: key, Value: value})
	}

	return root, nil
}
