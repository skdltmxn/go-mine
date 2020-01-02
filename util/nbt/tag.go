package nbt

import (
	"errors"
	"reflect"
)

func getTargetFieldNames(t reflect.Type) (map[string]int, error) {
	if t.Kind() != reflect.Struct {
		return nil, errors.New("nbt: struct must be given")
	}

	n := t.NumField()
	names := make(map[string]int)

	for i := 0; i < n; i++ {
		f := t.Field(i)
		tag := f.Tag.Get("nbt")
		if len(tag) == 0 {
			names[f.Name] = i
		} else {
			names[tag] = i
		}
	}

	return names, nil
}
