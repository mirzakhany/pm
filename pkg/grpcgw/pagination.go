package grpcgw

import "github.com/mirzakhany/pm/pkg/config"

var (
	maxPerPage = config.RegisterInt64("api.max_per_page", 100)
	minPerPage = config.RegisterInt64("api.min_per_page", 1)
	perPage    = config.RegisterInt64("api.per_page", 10)
)

// GetOffsetAndLimit return the offset and limit variable from the request, if not available
// return the default value
func GetOffsetAndLimit(offset, limit int64) (int64, int64) {

	if limit > maxPerPage.Int64() || limit < minPerPage.Int64() {
		limit = perPage.Int64()
	}
	return offset, limit
}
