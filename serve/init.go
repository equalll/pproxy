package serve
import "github.com/equalll/mydebug"

// 系统版本
var PproxyVersion string

func init() {mydebug.INFO()
	PproxyVersion = GetVersion()
}
