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

	DescribeTable("legal operators can be converted to expressions",

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
	)

	Context("legal filters", func() {
		movie := &model.Movie{
			Kind:  model.ShowKindMovie,
			Title: "Back to the Future",
		}

		FDescribeTable("can be evaluated as expressions",
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

			// Entry("multi-field nested filter",
			// 	`[{"kind": {"eq": "SERIES"}}, {"or": [{"kind": {"eq": "MOVIE"}}, {"and": [{"title": {"eq": "Back to the Future"}}]}]}]`,
			// 	`(kind == "SERIES") || ((kind == "MOVIE") && (title == "Back to the Future"))`,
			// 	movie,
			// ),
		)
	})

	//
	//			filter:     [kind: {eq: "SERIES"}, or: [kind: {eq: "MOVIE"}, and: {title: {eq: "Back to the Future"}}]]
	//		  expression: (kind == "SERIES") || ((kind == "MOVIE") && (title == "Back to the Future"))
})
