package internal

import (
	"fmt"
	"runtime"
)

var (
	mayor = "0"
	minor = "0"
	patch = "0"
	build = ""
)

func GetVersion() string {
	return fmt.Sprintf("%s.%s.%s build=%s os=%s arch=%s", mayor, minor, patch, build, runtime.GOARCH, runtime.GOOS)
}
