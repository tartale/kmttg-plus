package model

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/tartale/go/pkg/structs"
)

const maxExpressionDepth = 10

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

func (f *FilterOperator) Filter() (bool, error) {

	return false, nil
}

// Example:
//
//	f:          {eq: "foo"}
//	expression: `== "foo"`
func (f *FilterOperator) AsExpression() (expression string, err error) {

	operatorJsonBytes, err := json.Marshal(f)
	if err != nil {
		return
	}
	expression = format(string(operatorJsonBytes))
	if f.And != nil {
		expression = "&&"
	}
	if f.Or != nil {
		expression = "||"
	}

	return
}

// Example:
//
//		s:          {kind: {eq: "SERIES"}}
//	  expression: `kind == "SERIES"`
func (s *ShowFilter) AsExpression() (string, error) {

	var expressions []string
	filterWalkFn := func(field reflect.StructField, value reflect.Value) error {

		if value.IsNil() {
			return nil
		}

		switch val := value.Interface().(type) {
		case *FilterOperator:
			operatorExpression, err := val.AsExpression()
			if err != nil {
				return err
			}
			fieldName := jsonNameForReflectField(field)
			expressions = append(expressions, fmt.Sprintf("%s %s", fieldName, operatorExpression))
		}

		return nil
	}
	structs.Walk(s, filterWalkFn)

	return strings.Join(expressions, " && "), nil
}

// Example:
//
//		s:          {kind: {eq: "SERIES"}}
//		show:       {kind: "MOVIE", title: "Back to the Future"}
//	  vars:       {kind => "MOVIE"}
//	                    ^^ title is not in the map, since it's not in the filter
func (s *ShowFilter) ExtractVariables(input Show) (map[string]any, error) {

	vars := map[string]any{}
	filterWalkFn := func(filterField reflect.StructField, filterValue reflect.Value) error {

		if filterValue.IsNil() {
			return nil
		}
		switch filterValue.Interface().(type) {
		case *FilterOperator:
			showField, ok := structs.New(input).FieldOk(filterField.Name)
			if !ok {
				panic(fmt.Errorf("filter contains a field that is not in the input: %s", filterField.Name))
			}
			showFieldName := jsonNameForStructsField(showField)
			vars[showFieldName] = showField.Value()
		}

		return nil
	}
	structs.Walk(s, filterWalkFn)

	return vars, nil
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

type TivoFilterFn = func(t *Tivo) bool

func NewTivoFilter(f *TivoFilter) TivoFilterFn {

	return func(t *Tivo) bool {

		if f == nil {
			return true
		}
		if f.Name != nil {
			return true
		}

		return false
	}
}
