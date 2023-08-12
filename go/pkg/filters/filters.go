package filters

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/tartale/go/pkg/structs"
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

var (
	whitespaces      = regexp.MustCompile(`\s+`)
	operatorReplacer = strings.NewReplacer(
		`{`, ` `,
		`}`, ` `,
		`:`, ` `,
		`"eq"`, `==`,
		`"ne"`, `!=`,
		`"lt"`, `<`,
		`"gt"`, `>`,
		`"lte"`, `<=`,
		`"gte"`, `>=`,
		`"matches"`, `=~`,
		`"and"`, `&&`,
		`"or"`, `||`,
	)
)

func Filter[T any](slice []T, f func(T) bool) []T {

	var result []T
	for _, item := range slice {
		if f(item) {
			result = append(result, item)
		}
	}

	return result
}

// Example:
//
//	operator:   {eq: "foo"}
//	expression: `== "foo"`
func OperatorExpression(operator *model.FilterOperator) (expression string, err error) {

	operatorJsonBytes, err := json.Marshal(operator)
	if err != nil {
		return
	}
	expression = format(string(operatorJsonBytes))
	if operator.And != nil {
		expression = "&&"
	}
	if operator.Or != nil {
		expression = "||"
	}

	return
}

// Example:
//
//		filter:     {kind: {eq: "SERIES"}}
//	  expression: `kind == "SERIES"`
func FilterExpression(filter any) (string, error) {

	var expressions []string
	filterWalkFn := func(field reflect.StructField, value reflect.Value) error {

		if value.IsNil() {
			return nil
		}

		switch val := value.Interface().(type) {
		case *model.FilterOperator:
			operatorExpression, err := OperatorExpression(val)
			if err != nil {
				return err
			}
			fieldName := jsonNameForReflectField(field)
			expressions = append(expressions, fmt.Sprintf("%s %s", fieldName, operatorExpression))
		}

		return nil
	}
	structs.Walk(filter, filterWalkFn)

	return strings.Join(expressions, " && "), nil
}

// Example:
//
//		filter:     {kind: {eq: "SERIES"}}
//		input:      {kind: "MOVIE", title: "Back to the Future"}
//	  values:     {kind => "MOVIE"}
//	                    ^^ title is not in the map, since it's not in the filter
func GetValues(filter, input any) (map[string]any, error) {

	values := map[string]any{}
	filterWalkFn := func(filterField reflect.StructField, filterValue reflect.Value) error {

		if filterValue.IsNil() {
			return nil
		}
		switch filterValue.Interface().(type) {
		case *model.FilterOperator:
			showField, ok := structs.New(input).FieldOk(filterField.Name)
			if !ok {
				panic(fmt.Errorf("filter contains a field that is not in the input: %s", filterField.Name))
			}
			showFieldName := jsonNameForStructsField(showField)
			values[showFieldName] = showField.Value()
		}

		return nil
	}
	structs.Walk(filter, filterWalkFn)

	return values, nil
}

func format(expression string) string {

	expression = operatorReplacer.Replace(expression)
	expression = whitespaces.ReplaceAllString(expression, " ")
	expression = strings.Trim(expression, " ")

	return expression
}

func jsonNameForReflectField(field reflect.StructField) string {
	jsonTag, ok := field.Tag.Lookup("json")

	if !ok {
		panic(fmt.Errorf("missing json tag on field: %s", field.Name))
	}

	return strings.Split(jsonTag, ",")[0]
}

func jsonNameForStructsField(field *structs.Field) string {

	jsonTag := field.Tag("json")
	if jsonTag == "" {
		panic(fmt.Errorf("missing json tag on field: %s", field.Name()))
	}

	return strings.Split(jsonTag, ",")[0]
}
