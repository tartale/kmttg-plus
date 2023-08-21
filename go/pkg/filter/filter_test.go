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

	Context("legal filters", func() {
		movie := &model.Movie{
			Kind:  model.ShowKindMovie,
			Title: "Back to the Future",
		}

		DescribeTable("can be evaluated as expressions",
			func(showFiltersJson string, show model.Show) {

				var showFilters []*model.ShowFilter
				err := json.Unmarshal([]byte(showFiltersJson), &showFilters)
				Expect(err).ToNot(HaveOccurred())

				expression := filter.GetExpression(showFilters)
				values := filter.GetValues(showFilters, show)
				eval, err := gval.Evaluate(expression, values)

				Expect(err).ToNot(HaveOccurred())
				Expect(eval.(bool)).To(Equal(true))
			},

			Entry("simple filter",
				`[{"kind": {"eq": "MOVIE"}}]`,
				movie,
			),

			Entry("multi-field filter with implied logical 'and'",
				`[{"kind": {"eq": "MOVIE"}}, {"title": {"matches": "Back to the .*"}}]`,
				movie,
			),

			Entry("multi-field filter with explicit logical 'or'",
				`[{"kind": {"eq": "SERIES"}}, {"or": [{"title": {"matches": "Back to the .*"}}]}]`,
				movie,
			),

			Entry("multi-field nested filter",
				`[{"kind": {"eq": "SERIES"}}, {"or": [{"kind": {"eq": "MOVIE"}}, {"and": [{"title": {"eq": "Back to the Future"}}]}]}]`,
				movie,
			),
		)
	})
})
