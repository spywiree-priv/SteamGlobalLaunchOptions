package vdf

import (
	"iter"
	"maps"
	"strings"

	"github.com/guregu/null/v5"
)

type KeyValue struct {
	Key      string      `json:"key"`
	Value    null.String `json:"value"`
	Children []KeyValue  `json:"children"`
}

func (kv *KeyValue) GetChild(key string) (*KeyValue, bool) {
	for child := range kv.ChildrenIter() {
		if child.Key == key {
			return child, true
		}
	}
	return nil, false
}

func (kv *KeyValue) GetChildren(key string) []*KeyValue {
	var children []*KeyValue
	for child := range kv.ChildrenIter() {
		if child.Key == key {
			children = append(children, child)
		}
	}
	return children
}

type ErrKeyNotFound []string

func (err ErrKeyNotFound) Error() string {
	return "key: " + strings.Join(err, " - ") + " not found"
}

func (kv *KeyValue) GetChildByPath(path ...string) (*KeyValue, error) {
	child := kv
	var found bool
	for i, key := range path {
		child, found = child.GetChild(key)
		if !found {
			return nil, ErrKeyNotFound(path[:i+1])
		}
	}
	return child, nil
}

func (kv *KeyValue) SetChild(value KeyValue) {
	for child := range kv.ChildrenIter() {
		if child.Key == value.Key {
			*child = value
			return
		}
	}
	kv.Children = append(kv.Children, value)
}

func (kv *KeyValue) HasChild(key string) bool {
	for child := range kv.ChildrenIter() {
		if child.Key == key {
			return true
		}
	}
	return false
}

func (kv *KeyValue) HasNonZeroChild(key string) bool {
	for child := range kv.ChildrenIter() {
		if child.Key == key && child.Value.ValueOrZero() != "" {
			return true
		}
	}
	return false
}

func (kv *KeyValue) ChildrenIter() iter.Seq[*KeyValue] {
	return func(yield func(*KeyValue) bool) {
		for i := range kv.Children {
			if !yield(&kv.Children[i]) {
				return
			}
		}
	}
}

func (kv *KeyValue) ToMap() map[string]any {
	if kv.Value.Valid {
		return map[string]any{kv.Key: kv.Value.String}
	}

	m := make(map[string]any)
	for _, child := range kv.Children {
		maps.Copy(m, child.ToMap())
	}
	return map[string]any{kv.Key: m}
}
