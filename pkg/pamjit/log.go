package pamjit

import (
	"fmt"
	"strings"
)

func PrintLine(leftPadding int, format string, args ...any) {
	fmt.Printf("%s%s\n", strings.Repeat(" ", leftPadding), fmt.Sprintf(format, args...))
}
