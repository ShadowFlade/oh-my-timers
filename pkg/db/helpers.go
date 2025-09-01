package db

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

type Helper struct {
}

func (this *Helper) GenerateInsertQuery(tableName string, s interface{}) (string, error) {
	value := reflect.ValueOf(s)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return "", fmt.Errorf("expected a struct, got %T", s)
	}

	var columns []string
	sb := strings.Builder{}

	for i := 0; i < value.NumField(); i++ {
		fieldValue := fmt.Sprint(value.Field(i))
		tag := value.Type().Field(i).Tag.Get("db")

		if tag == "" || tag == "-" {
			continue // скипаем не db тэги
		}

		columns = append(columns, tag)

		log.Print(fieldValue," FIELD VALUE")

		if fieldValue == "" {
			sb.WriteString(" ,")
		} else {
			sb.WriteString(fieldValue + ",")
		}
	}

	if len(columns) == 0 {
		return "", fmt.Errorf("no fields with `db` tags found")
	}
	values := sb.String()

	if string(values[len(values)-1:]) == "," {
		values = values[:len(values)-1]
	}
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s);",
		tableName,
		strings.Join(columns, ", "),
		values,
	)
	return query, nil
}
