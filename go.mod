module github.com/rcrowley/feed

go 1.25.4

require (
	github.com/rcrowley/mergician v0.0.0-20251130073118-e00557126233
	golang.org/x/net v0.47.0
)

require (
	github.com/yuin/goldmark v1.7.13 // indirect
	golang.org/x/exp v0.0.0-20251125195548-87e1e737ad39 // indirect
)

// TEMPORARY: this -n scheduling feature depends on files.PublishStatus,
// which lives on the editorial-workflow branch of mergician
// (rcrowley/mergician PR). Remove this replace once that change is merged
// and the require above is bumped to a version that includes it.
replace github.com/rcrowley/mergician => github.com/rbotley/mergician v0.0.0-20260605061827-8cc9af92e627
