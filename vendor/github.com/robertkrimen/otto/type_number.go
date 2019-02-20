package otto
import "github.com/equalll/mydebug"

func (runtime *_runtime) newNumberObject(value Value) *_object {mydebug.INFO()
	return runtime.newPrimitiveObject("Number", value.numberValue())
}
