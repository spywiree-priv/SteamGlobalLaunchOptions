package vdf

import (
	"io"
	"strings"
)

func (kv *KeyValue) WriteIndent(w io.Writer, level int) error {
	indent := strings.Repeat("\t", level)

	if _, err := io.WriteString(w, indent); err != nil {
		return err
	}
	if _, err := io.WriteString(w, Quote(kv.Key)); err != nil {
		return err
	}
	if kv.Value.Valid {
		if _, err := io.WriteString(w, "\t\t"); err != nil {
			return err
		}
		if _, err := io.WriteString(w, Quote(kv.Value.String)); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(w, "\n"); err != nil {
		return err
	}

	if len(kv.Children) > 0 {
		if _, err := io.WriteString(w, indent); err != nil {
			return err
		}
		if _, err := io.WriteString(w, "{\n"); err != nil {
			return err
		}
		for _, child := range kv.Children {
			if err := child.WriteIndent(w, level+1); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(w, indent); err != nil {
			return err
		}
		if _, err := io.WriteString(w, "}\n"); err != nil {
			return err
		}
	}

	return nil
}

func (kv *KeyValue) Write(w io.Writer) error {
	return kv.WriteIndent(w, 0)
}

func (kv *KeyValue) String() string {
	var sb strings.Builder
	_ = kv.Write(&sb)
	return sb.String()
}
