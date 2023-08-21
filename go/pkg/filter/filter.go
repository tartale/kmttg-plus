package filter

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/tartale/go/pkg/structs"
	"golang.org/x/exp/maps"
)

var (
	whitespace        = regexp.MustCompile(`\s+`)
	quotedFields      = regexp.MustCompile(`"(\w+)":`)
	doubleLeftParens  = regexp.MustCompile(`\( \(`)
	doubleRightParens = regexp.MustCompile(`\) \)`)
	replacer          = strings.NewReplacer(
		`eq`, ` == `,
		`ne`, ` != `,
		`lte`, ` <= `,
		`gte`, ` >= `,
		`lt`, ` < `,
		`gt`, ` > `,
		`matches`, ` =~ `,
		`,{or`, ` ) || ( `,
		`,{and`, ` ) && ( `,
		`,`, ` ) && ( `,
		`[`, `(`,
		`]`, `)`,
		`{`, ` `,
		`}`, ` `,
	)
	stringType = reflect.TypeOf("")
)

type Operator struct {
	Eq      any `json:"eq,omitempty"`
	Ne      any `json:"ne,omitempty"`
	Lt      any `json:"lt,omitempty"`
	Gt      any `json:"gt,omitempty"`
	Lte     any `json:"lte,omitempty"`
	Gte     any `json:"gte,omitempty"`
	Matches any `json:"matches,omitempty"`
}

// Examples:
//
//	operator:   {eq: "foo"}
//	expression: == "foo"
func OperatorExpression(operator *Operator) (expression string) {

	operatorJsonBytes, err := json.Marshal(operator)
	if err != nil {
		panic(fmt.Errorf("unexpected error when marshaling operator: %w", err))
	}
	expression = format(string(operatorJsonBytes))

	return
}

// Examples:
//
//			filter:     {kind: {eq: "SERIES"}}
//		  expression: (kind == "SERIES")
//
//			filter:     [kind: {eq: "SERIES"}, title: {eq: "Back to the Future"}]
//		  expression: (kind == "SERIES") && (title == "Back to the Future")
//		                    ^^ when multiple fields are given without a logical operator,
//	                         the default logical operator is "and"
//
//			filter:     [kind: {eq: "SERIES"}, or: {kind: {eq: "EPISODE"}]
//		  expression: (kind == "SERIES") || (kind == "EPISODE")
//
//			filter:     [kind: {eq: "SERIES"}, or: [kind: {eq: "MOVIE"}, and: {title: {eq: "Back to the Future"}}]]
//		  expression: (kind == "SERIES") || ((kind == "MOVIE") && (title == "Back to the Future"))
func GetExpression(filter any) string {

	filterValue := reflect.ValueOf(filter)
	if !structs.IsSlice(filterValue) {
		filter = []any{filter}
	}
	filterBytes, err := json.Marshal(filter)
	if err != nil {
		panic(fmt.Errorf("unexpected error when marshaling filter: %w", err))
	}

	filterJson := string(filterBytes)
	expression := format(filterJson)

	return expression
}

// Example:
//
//		filter:     {kind: {eq: "SERIES"}}
//		input:      {kind: "MOVIE", title: "Back to the Future"}
//	  values:     {kind => "MOVIE"}
//	                    ^^ title is not in the map, since it's not in the filter
func GetValues(filter, input any) map[string]any {

	filterValue := reflect.ValueOf(filter)
	if !structs.IsSlice(filterValue) {
		filter = []any{filter}
	}

	values := map[string]any{}
	for i := 0; i < filterValue.Len(); i++ {
		f := filterValue.Index(i).Interface()
		v := getValues(f, input)
		maps.Copy(values, v)
	}

	return values
}

func getValues(filter, input any) map[string]any {

	values := map[string]any{}
	filterWalkFn := func(filterField reflect.StructField, filterValue reflect.Value) error {

		if filterValue.IsNil() {
			return nil
		}
		switch filterValue.Interface().(type) {
		case *Operator:

			inputField, ok := structs.New(input).FieldOk(filterField.Name)
			if !ok {
				panic(fmt.Errorf("filter contains a field that is not in the input: %s", filterField.Name))
			}
			inputFieldName := inputField.TagRoot("json")
			inputFieldValue := inputField.Value()
			inputFieldReflectValue := reflect.ValueOf(inputFieldValue)
			if inputFieldReflectValue.CanConvert(stringType) {
				inputFieldValue = inputFieldReflectValue.Convert(stringType).Interface()
			}
			values[inputFieldName] = inputFieldValue
		}

		return nil
	}

	structs.Walk(filter, filterWalkFn)

	return values
}

func format(expression string) string {

	expression = quotedFields.ReplaceAllString(expression, "$1")
	expression = replacer.Replace(expression)
	expression = whitespace.ReplaceAllString(expression, " ")
	expression = doubleLeftParens.ReplaceAllString(expression, "(")
	expression = doubleRightParens.ReplaceAllString(expression, ")")
	expression = strings.Trim(expression, " ")

	return expression
}
