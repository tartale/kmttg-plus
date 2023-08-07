package model

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"go.uber.org/zap"
)

var _ = Describe("Show Helpers", func() {

	It("can merge episodes from a list of recordings into a series", func() {
		unmergedShows := []Show{
			&Episode{
				Kind:         ShowKindEpisode,
				Title:        "Law and Order",
				EpisodeTitle: "Bad Guys",
			},
			&Episode{
				Kind:         ShowKindEpisode,
				Title:        "CSI: Miami",
				EpisodeTitle: "More Bad Guys",
			},
			&Movie{
				Kind:        ShowKindMovie,
				Title:       "Back to the Future",
				Description: "",
				MovieYear:   1985,
			},
			&Episode{
				Kind:         ShowKindEpisode,
				Title:        "Law and Order",
				EpisodeTitle: "Even More Bad Guys",
			},
		}
		mergedShows := MergeEpisodes(unmergedShows)
		Expect(mergedShows).To(HaveLen(3))
		logz.Logger.Info("", zap.Any("mergedShows", mergedShows))
	})
})
