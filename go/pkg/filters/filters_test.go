package filters

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

var _ = Describe("Filters", func() {

	DescribeTable("legal operators can be converted to expressions",
		func(operator model.FilterOperator, expectedExpression string) {

			expression, err := OperatorExpression(&operator)

			Expect(err).ToNot(HaveOccurred())
			Expect(expression).To(Equal(expectedExpression))
		},

		Entry("equals string", model.FilterOperator{Eq: "foo"}, `== "foo"`),
		Entry("not equals string", model.FilterOperator{Ne: "foo"}, `!= "foo"`),
		Entry("less than string", model.FilterOperator{Lt: "foo"}, `< "foo"`),
		Entry("greater than string", model.FilterOperator{Gt: "foo"}, `> "foo"`),
		Entry("less than or equal to string", model.FilterOperator{Lte: "foo"}, `<= "foo"`),
		Entry("greater than or equal to string", model.FilterOperator{Gte: "foo"}, `>= "foo"`),
		Entry("matches string", model.FilterOperator{Matches: ".*foo.*"}, `=~ ".*foo.*"`),

		Entry("equals integer", model.FilterOperator{Eq: 10}, `== 10`),
		Entry("not equals integer", model.FilterOperator{Ne: 10}, `!= 10`),
		Entry("less than integer", model.FilterOperator{Lt: 10}, `< 10`),
		Entry("greater than integer", model.FilterOperator{Gt: 10}, `> 10`),
		Entry("less than or equal to integer", model.FilterOperator{Lte: 10}, `<= 10`),
		Entry("greater than or equal to integer", model.FilterOperator{Gte: 10}, `>= 10`),

		Entry("equals float", model.FilterOperator{Eq: 10.9}, `== 10.9`),
		Entry("not equals float", model.FilterOperator{Ne: 10.9}, `!= 10.9`),
		Entry("less than float", model.FilterOperator{Lt: 10.9}, `< 10.9`),
		Entry("greater than float", model.FilterOperator{Gt: 10.9}, `> 10.9`),
		Entry("less than or equal to float", model.FilterOperator{Lte: 10.9}, `<= 10.9`),
		Entry("greater than or equal to float", model.FilterOperator{Gte: 10.9}, `>= 10.9`),

		Entry("and", model.FilterOperator{And: &model.FilterOperator{}}, `&&`),
		Entry("or", model.FilterOperator{Or: &model.FilterOperator{}}, `||`),
	)

	Context("simple filter", func() {

		showFilter := &model.ShowFilter{
			Kind: &model.FilterOperator{
				Eq: model.ShowKindMovie,
			},
		}

		show := model.Movie{
			Kind:  model.ShowKindMovie,
			Title: "Back to the Future",
		}

		It("can be converted to an expression", func() {

			expression, err := FilterExpression(showFilter)

			Expect(err).ToNot(HaveOccurred())
			Expect(expression).To(Equal(`kind == "MOVIE"`))
		})

		It("can get variables from an input value", func() {

			vars, err := GetValues(showFilter, show)

			Expect(err).ToNot(HaveOccurred())
			Expect(vars).To(HaveKeyWithValue("kind", model.ShowKindMovie))
			Expect(vars).ToNot(HaveKeyWithValue("title", model.ShowKindMovie))
			Expect(vars).ToNot(HaveKey("movieYear"))
		})
	})

	Context("multiple field filter", func() {

		showFilter := &model.ShowFilter{
			Kind: &model.FilterOperator{
				Eq: model.ShowKindSeries,
			},
			Title: &model.FilterOperator{
				Matches: "Back to the .*",
			},
		}

		show := model.Movie{
			Kind:        model.ShowKindMovie,
			Title:       "Back to the Future",
			Description: "Doc and Marty hijinks",
		}

		It("can be converted to an expression", func() {

			expression, err := FilterExpression(showFilter)

			Expect(err).ToNot(HaveOccurred())
			Expect(expression).To(Equal(`kind == "SERIES" && title =~ "Back to the .*"`))
		})

		It("can get variables from an input value", func() {

			vars, err := GetValues(showFilter, show)

			Expect(err).ToNot(HaveOccurred())
			Expect(vars).To(HaveKeyWithValue("kind", model.ShowKindMovie))
			Expect(vars).To(HaveKeyWithValue("title", "Back to the Future"))
			Expect(vars).ToNot(HaveKey("movieYear"))
		})
	})
})
