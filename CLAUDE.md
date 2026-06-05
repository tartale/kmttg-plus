# Claude Code Notes

## go/pkg/shows — refactoring complete

The `shows` package wraps model types (`model.Movie`, `model.Series`, `model.Episode`) to carry
TiVo-specific metadata (`Details`: raw recording/collection data, object ID, Tivo connection) without
exposing it through the GraphQL API.

**Current design (Option B — single `DetailedShow` wrapper):**

```go
type DetailedShow struct {
    model.Show                        // embedded interface — delegates all Show methods
    Details        Details            // TiVo metadata for this show
    EpisodeDetails map[string]Details // per-episode metadata, keyed by episode ID; series only
}
```

- `AsApiType(show)` unwraps to `ds.Show` (the plain model type for API responses)
- `GetDetails(show)` returns `&ds.Details` via a single type assertion, no switch
- `NewFilterFn` unwraps `*DetailedShow` to its inner `model.Show` before passing to the
  reflection-based filter package (which cannot see fields through an embedded interface)
- `MarshalJSON` flattens the inner show's fields alongside `details` and `episodeDetails`
  into one JSON object; `UnmarshalShowFromJSON` reconstructs the wrapper from that format
- Series episodes: `model.Series.Episodes []*model.Episode` carries the API slice;
  `EpisodeDetails` carries the corresponding per-episode TiVo metadata
- `GetEpisodesForSeries` reconstructs per-episode `*DetailedShow` values on demand from
  the series's `Episodes` slice and its `EpisodeDetails` map

**Test coverage:** `go/pkg/shows/shows_test.go` covers all public functions with unit tests
and two integration tests (`TestIntegration_*`) that require `KMTTG_TEST_TIVO` to be set.
