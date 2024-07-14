package passwd

import (
	"strings"

	uuid "github.com/iris-contrib/go.uuid"
)

func UUID() string {
	u, _ := uuid.NewV4()
	return strings.ReplaceAll(u.String(), "-", "")
}