# Claude Code Notes

## go/pkg/shows — refactoring in progress

The `shows` package wraps model types (`model.Movie`, `model.Series`, `model.Episode`) to carry
TiVo-specific metadata (`Details`: raw recording/collection data, object ID, Tivo connection) without
exposing it through the GraphQL API. The current approach defines three parallel wrapper structs
(`movie`, `series`, `episode`) that embed the model type plus a `Details` field.

**Pain points with the current approach:**
- Every function (`Clone`, `MarshalShowToJSON`, `UnmarshalShowFromJSON`, `AsApiType`, `WithImageURL`)
  repeats a 3-case type switch on `GetKind()`
- `series` has dual `Episodes` fields: `*model.Series.Episodes []*model.Episode` (API) and
  `series.Episodes []*episode` (internal), requiring synchronization on conversion
- `AsApiType()` exists only to strip the wrapper before returning data to the API

**Options under consideration:**

### Option A — Side-channel details store
Drop wrapper types entirely. Return pure `model.Movie/Series/Episode` from `shows.New()` and
register their `Details` in a package-level `sync.Map` keyed by show ID. `AsApiType()` disappears,
the dual-Episodes problem goes away, no type switching in shows package.
- Tradeoff: cache serialization needs rethinking (details currently live alongside model JSON);
  store needs explicit lifecycle management (register on load, delete on cache clear).

### Option B — Single `DetailedShow` wrapper (embedding model.Show interface)
Replace three wrapper structs with one:
```go
type DetailedShow struct {
    model.Show  // embedded interface — delegates all Show methods automatically
    Details Details
}
```
`DetailedShow` satisfies `model.Show` automatically. `AsApiType()` becomes just `ds.Show`.
- Tradeoff: JSON marshaling needs custom `MarshalJSON`/`UnmarshalJSON` to flatten the interface value.
- Tradeoff: series episodes — `model.Series.Episodes` is `[]*model.Episode` with no Details;
  per-episode Details need a separate field or a flat (non-nested) episode store.
- Type assertions to concrete types are still needed downstream.

### Option C — Generic wrapper (least invasive)
```go
type Annotated[T any] struct { Inner T; Details Details }
```
Collapses three structs into one generic type; type-switch boilerplate reduces to generic helpers.
Can't embed a type parameter in Go so callers use `.Inner` instead of direct field access.
Cache serialization stays the same. Less conceptual cleanup than A or B.

**Leading candidate:** Option B (`DetailedShow` embedding `model.Show`) with flat episode handling
(episodes stored as top-level `DetailedShow` entries; `Series.Episodes` populated from those at
API response time rather than carried inside the wrapper).
