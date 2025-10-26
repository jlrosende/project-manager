package internal

import (
	"fmt"
	"runtime"
)

var (
	version = ""
	commit  = ""
	date    = ""
	builtBy = ""
)

func GetVersion() string {
	return fmt.Sprintf(
		"%s (%s) [by=%s os=%s arch=%s date=%s]",
		version,
		commit,
		builtBy,
		runtime.GOARCH,
		runtime.GOOS,
		date,
	)
}
