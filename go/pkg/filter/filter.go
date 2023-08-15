package filter

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/tartale/go/pkg/structs"
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

type Operator struct {
	Eq      interface{} `json:"eq,omitempty"`
	Ne      interface{} `json:"ne,omitempty"`
	Lt      interface{} `json:"lt,omitempty"`
	Gt      interface{} `json:"gt,omitempty"`
	Lte     interface{} `json:"lte,omitempty"`
	Gte     interface{} `json:"gte,omitempty"`
	Matches interface{} `json:"matches,omitempty"`
	And     *Operator   `json:"-"`
	Or      *Operator   `json:"-"`
}

// Example:
//
//	operator:   {eq: "foo"}
//	expression: `== "foo"`
func OperatorExpression(operator *Operator) (expression string) {

	operatorJsonBytes, err := json.Marshal(operator)
	if err != nil {
		panic(fmt.Errorf("unexpected error when marshaling operator: %w", err))
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
func GetExpression(filter any) string {

	var expressions []string
	filterWalkFn := func(field reflect.StructField, value reflect.Value) error {

		if value.IsNil() {
			return nil
		}

		switch val := value.Interface().(type) {
		case *Operator:
			operatorExpression := OperatorExpression(val)
			fieldName := jsonNameForReflectField(field)
			expressions = append(expressions, fmt.Sprintf("%s %s", fieldName, operatorExpression))
		}

		return nil
	}
	structs.Walk(filter, filterWalkFn)

	return strings.Join(expressions, " && ")
}

// Example:
//
//		filter:     {kind: {eq: "SERIES"}}
//		input:      {kind: "MOVIE", title: "Back to the Future"}
//	  values:     {kind => "MOVIE"}
//	                    ^^ title is not in the map, since it's not in the filter
func GetValues(filter, input any) map[string]any {

	values := map[string]any{}
	filterWalkFn := func(filterField reflect.StructField, filterValue reflect.Value) error {

		if filterValue.IsNil() {
			return nil
		}
		switch filterValue.Interface().(type) {
		case *Operator:
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

	return values
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
