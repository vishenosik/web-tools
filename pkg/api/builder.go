package api

import (
	"fmt"
	"path"
)

const (
	prefix string = "/api"
)

func ApiV1(routeParts ...string) string {
	return buildApi(1, routeParts...)
}

func buildApi(version uint8, routeParts ...string) string {
	ver := fmt.Sprintf("v%v", version)
	route := path.Join(prefix, ver)
	for i := range routeParts {
		route = path.Join(route, routeParts[i])
	}
	return route
}
