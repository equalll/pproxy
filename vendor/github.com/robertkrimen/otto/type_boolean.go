package otto
import "github.com/equalll/mydebug"

import (
	"strconv"
)

func (runtime *_runtime) newBooleanObject(value Value) *_object {mydebug.INFO()
	return runtime.newPrimitiveObject("Boolean", toValue_bool(value.bool()))
}

func booleanToString(value bool) string {mydebug.INFO()
	return strconv.FormatBool(value)
}
