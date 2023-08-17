package filter_test

import (
	"encoding/json"

	"github.com/PaesslerAG/gval"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tartale/kmttg-plus/go/pkg/filter"
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

var _ = Describe("Filter", func() {

	PDescribeTable("legal operators can be converted to expressions",

		func(operator *filter.Operator, expectedExpression string) {

			expression := filter.OperatorExpression(operator)
			Expect(expression).To(Equal(expectedExpression))
		},

		Entry("equals string", &filter.Operator{Eq: "foo"}, `== "foo"`),
		Entry("not equals string", &filter.Operator{Ne: "foo"}, `!= "foo"`),
		Entry("less than string", &filter.Operator{Lt: "foo"}, `< "foo"`),
		Entry("greater than string", &filter.Operator{Gt: "foo"}, `> "foo"`),
		Entry("less than or equal to string", &filter.Operator{Lte: "foo"}, `<= "foo"`),
		Entry("greater than or equal to string", &filter.Operator{Gte: "foo"}, `>= "foo"`),
		Entry("matches string", &filter.Operator{Matches: ".*foo.*"}, `=~ ".*foo.*"`),

		Entry("equals integer", &filter.Operator{Eq: 10}, `== 10`),
		Entry("not equals integer", &filter.Operator{Ne: 10}, `!= 10`),
		Entry("less than integer", &filter.Operator{Lt: 10}, `< 10`),
		Entry("greater than integer", &filter.Operator{Gt: 10}, `> 10`),
		Entry("less than or equal to integer", &filter.Operator{Lte: 10}, `<= 10`),
		Entry("greater than or equal to integer", &filter.Operator{Gte: 10}, `>= 10`),

		Entry("equals float", &filter.Operator{Eq: 10.9}, `== 10.9`),
		Entry("not equals float", &filter.Operator{Ne: 10.9}, `!= 10.9`),
		Entry("less than float", &filter.Operator{Lt: 10.9}, `< 10.9`),
		Entry("greater than float", &filter.Operator{Gt: 10.9}, `> 10.9`),
		Entry("less than or equal to float", &filter.Operator{Lte: 10.9}, `<= 10.9`),
		Entry("greater than or equal to float", &filter.Operator{Gte: 10.9}, `>= 10.9`),

		Entry("or", &filter.Operator{Or: &filter.Operator{}}, `||`),
	)

	Context("legal filters", func() {
		movie := &model.Movie{
			Kind:  model.ShowKindMovie,
			Title: "Back to the Future",
		}

		DescribeTable("legal operators can be converted to expressions",
			func(showFiltersJson, expectedExpression string, show model.Show) {

				var showFilters []*model.ShowFilter
				err := json.Unmarshal([]byte(showFiltersJson), &showFilters)
				Expect(err).ToNot(HaveOccurred())

				expression := filter.GetExpression(showFilters)
				Expect(expression).To(Equal(expectedExpression))

				values := filter.GetValues(showFilters, show)
				eval, err := gval.Evaluate(expression, values)

				Expect(err).ToNot(HaveOccurred())
				Expect(eval.(bool)).To(Equal(true))
			},

			Entry("simple filter",
				`[{"kind": {"eq": "MOVIE"}}]`,
				`( kind == "MOVIE" )`,
				movie,
			),

			Entry("multi-field filter with implied logical 'and'",
				`[{"kind": {"eq": "MOVIE"}}, {"title": {"matches": "Back to the .*"}}]`,
				`( kind == "MOVIE" ) && ( title =~ "Back to the .*" )`,
				movie,
			),

			Entry("multi-field filter with explicit logical 'or'",
				`[{"kind": {"eq": "SERIES"}}, {"or": [{"title": {"matches": "Back to the .*"}}]}]`,
				`( kind == "SERIES" ) || ( title =~ "Back to the .*" )`,
				movie,
			),
		)

	})
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

	// DescribeTable("legal operators can be converted to expressions",
	// 	func(showFilters []*model.ShowFilter, show model.Show, expectedExpression string) {

	// 		expression := filter.GetExpression(showFilters)
	// 		Expect(expression).To(Equal(expectedExpression))

	// 		values := filter.GetValues(showFilters, show)
	// 		eval, err := gval.Evaluate(expression, values)

	// 		Expect(err).ToNot(HaveOccurred())
	// 		Expect(eval.(bool)).To(Equal(true))

	// 	},

	// 	Entry("simple filter",
	// 		[]*model.ShowFilter{
	// 			{Kind: &filter.Operator{
	// 				Eq: model.ShowKindMovie,
	// 			},
	// 			}},
	// 		&model.Movie{
	// 			Kind:  model.ShowKindMovie,
	// 			Title: "Back to the Future",
	// 		},
	// 		`( kind == "MOVIE" )`,
	// 	),

	// 	Entry("multi-field filter with implied logical 'and'",
	// 		[]*model.ShowFilter{
	// 			{
	// 				Kind: &filter.Operator{
	// 					Eq: model.ShowKindMovie,
	// 				},
	// 			},
	// 			{
	// 				Title: &filter.Operator{
	// 					Matches: "Back to the .*",
	// 				},
	// 			},
	// 		},
	// 		&model.Movie{
	// 			Kind:        model.ShowKindMovie,
	// 			Title:       "Back to the Future",
	// 			Description: "Doc and Marty's hijinks",
	// 		},
	// 		`( kind == "MOVIE" ) && ( title =~ "Back to the .*" )`,
	// 	),

	// 	Entry("multi-field filter with explicit logical 'or'",
	// 		[]*model.ShowFilter{
	// 			{
	// 				Kind: &filter.Operator{
	// 					Eq: model.ShowKindSeries,
	// 				},
	// 			},
	// 			{
	// 				Or: &model.ShowFilter{
	// 					Title: &filter.Operator{
	// 						Matches: "Back to the .*",
	// 					},
	// 				},
	// 			},
	// 		},
	// 		&model.Movie{
	// 			Kind:        model.ShowKindMovie,
	// 			Title:       "Back to the Future",
	// 			Description: "Doc and Marty's hijinks",
	// 		},
	// 		`( kind == "SERIES" ) || ( title =~ "Back to the .*" )`,
	// 	),
	// 	Entry("nested filter",
	// 	),
	// )
})
