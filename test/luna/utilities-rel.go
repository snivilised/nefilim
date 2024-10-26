package luna

import (
	"strings"
)

// Yoke is similar to filepath.Join but it is meant specifically for relative file
// systems where the rules of a path are different; see fs.ValidPath
func Yoke(segments ...string) string {
	return strings.Join(segments, "/")
}
