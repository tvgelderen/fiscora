package repository

import "strings"

const DefaultFetchLimit = 25
const MaxFetchLimit = 1000

func NoRowsFound(err error) bool {
	return strings.Contains(err.Error(), "no rows in result set")
}
