package model

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filters", func() {

	DescribeTable("legal operators can be converted to expressions",
		func(operator FilterOperator, expectedExpression string) {

			expression, err := OperatorExpression(&operator)

			Expect(err).ToNot(HaveOccurred())
			Expect(expression).To(Equal(expectedExpression))
		},

		Entry("equals string", FilterOperator{Eq: "foo"}, `== "foo"`),
		Entry("not equals string", FilterOperator{Ne: "foo"}, `!= "foo"`),
		Entry("less than string", FilterOperator{Lt: "foo"}, `< "foo"`),
		Entry("greater than string", FilterOperator{Gt: "foo"}, `> "foo"`),
		Entry("less than or equal to string", FilterOperator{Lte: "foo"}, `<= "foo"`),
		Entry("greater than or equal to string", FilterOperator{Gte: "foo"}, `>= "foo"`),
		Entry("matches string", FilterOperator{Matches: ".*foo.*"}, `=~ ".*foo.*"`),

		Entry("equals integer", FilterOperator{Eq: 10}, `== 10`),
		Entry("not equals integer", FilterOperator{Ne: 10}, `!= 10`),
		Entry("less than integer", FilterOperator{Lt: 10}, `< 10`),
		Entry("greater than integer", FilterOperator{Gt: 10}, `> 10`),
		Entry("less than or equal to integer", FilterOperator{Lte: 10}, `<= 10`),
		Entry("greater than or equal to integer", FilterOperator{Gte: 10}, `>= 10`),

		Entry("equals float", FilterOperator{Eq: 10.9}, `== 10.9`),
		Entry("not equals float", FilterOperator{Ne: 10.9}, `!= 10.9`),
		Entry("less than float", FilterOperator{Lt: 10.9}, `< 10.9`),
		Entry("greater than float", FilterOperator{Gt: 10.9}, `> 10.9`),
		Entry("less than or equal to float", FilterOperator{Lte: 10.9}, `<= 10.9`),
		Entry("greater than or equal to float", FilterOperator{Gte: 10.9}, `>= 10.9`),

		Entry("and", FilterOperator{And: &FilterOperator{}}, `&&`),
		Entry("or", FilterOperator{Or: &FilterOperator{}}, `||`),
	)

	Context("simple filter", func() {

		showFilter := &ShowFilter{
			Kind: &FilterOperator{
				Eq: ShowKindMovie,
			},
		}

		show := Movie{
			Kind:  ShowKindMovie,
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
			Expect(vars).To(HaveKeyWithValue("kind", ShowKindMovie))
			Expect(vars).ToNot(HaveKeyWithValue("title", ShowKindMovie))
			Expect(vars).ToNot(HaveKey("movieYear"))
		})
	})

	Context("multiple field filter", func() {

		showFilter := &ShowFilter{
			Kind: &FilterOperator{
				Eq: ShowKindSeries,
			},
			Title: &FilterOperator{
				Matches: "Back to the .*",
			},
		}

		show := Movie{
			Kind:        ShowKindMovie,
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
			Expect(vars).To(HaveKeyWithValue("kind", ShowKindMovie))
			Expect(vars).To(HaveKeyWithValue("title", "Back to the Future"))
			Expect(vars).ToNot(HaveKey("movieYear"))
		})
	})
})
