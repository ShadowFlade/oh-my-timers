package db

import (
	"fmt"
	"reflect"
	"strings"
)

type Helpers struct {
}

func (this *Helpers) GenerateInsertQuery(table string, s interface{}) (string, error) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return "", fmt.Errorf("expected a struct, got %T", s)
	}

	var columns []string
	var values []string

	for i := 0; i < v.NumField(); i++ {
		value := fmt.Sprint(v.Field(i))
		tag := v.Type().Field(i).Tag.Get("db")

		if tag == "" || tag == "-" {
			continue // скипаем не db тэги
		}

		columns = append(columns, tag)
		values = append(values, value)
	}

	if len(columns) == 0 {
		return "", fmt.Errorf("no fields with `db` tags found")
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(columns, ", "),
		strings.Join(values, ", "),
	)
	return query, nil
}
