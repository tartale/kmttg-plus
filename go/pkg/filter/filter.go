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
	quotedFields = regexp.MustCompile(`"(\w+)":`)
	typeOfString = reflect.TypeOf("")
)

type Operator struct {
	Eq      any `json:"eq,omitempty"`
	Ne      any `json:"ne,omitempty"`
	Lte     any `json:"lte,omitempty"`
	Gte     any `json:"gte,omitempty"`
	Lt      any `json:"lt,omitempty"`
	Gt      any `json:"gt,omitempty"`
	Matches any `json:"matches,omitempty"`
}

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
			if inputFieldReflectValue.CanConvert(typeOfString) {
				inputFieldValue = inputFieldReflectValue.Convert(typeOfString).Interface()
			}
			values[inputFieldName] = inputFieldValue
		}

		return nil
	}

	structs.Walk(filter, filterWalkFn)

	return values
}

func removeQuotesOnFields(s string) string {

	return quotedFields.ReplaceAllString(s, "$1")
}

func replaceComparisonOperators(s string) string {

	s = regexp.MustCompile(`{eq(.*?)}`).ReplaceAllString(s, " == $1 ")
	s = regexp.MustCompile(`{ne(.*?)}`).ReplaceAllString(s, " != $1 ")
	s = regexp.MustCompile(`{lte(.*?)}`).ReplaceAllString(s, " <= $1 ")
	s = regexp.MustCompile(`{gte(.*?)}`).ReplaceAllString(s, " >= $1 ")
	s = regexp.MustCompile(`{lt(.*?)}`).ReplaceAllString(s, " < $1 ")
	s = regexp.MustCompile(`{gt(.*?)}`).ReplaceAllString(s, " > $1 ")
	s = regexp.MustCompile(`{matches(.*?)}`).ReplaceAllString(s, " =~ $1 ")

	return s
}

func replaceBrackets(s string) string {

	return strings.NewReplacer(
		`[`, `(`,
		`]`, `)`,
		`{`, `(`,
		`}`, `)`,
	).Replace(s)
}

func replaceLogicOperators(s string) string {

	return strings.NewReplacer(
		`,(or`, ` || (`,
		`,(and`, ` && (`,
		`,`, ` && `,
	).Replace(s)

}

func format(expression string) string {

	expression = removeQuotesOnFields(expression)
	expression = replaceComparisonOperators(expression)
	expression = replaceBrackets(expression)
	expression = replaceLogicOperators(expression)

	return expression
}
