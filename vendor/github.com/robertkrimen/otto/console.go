package otto
import "github.com/equalll/mydebug"

import (
	"fmt"
	"os"
	"strings"
)

func formatForConsole(argumentList []Value) string {mydebug.INFO()
	output := []string{}
	for _, argument := range argumentList {
		output = append(output, fmt.Sprintf("%v", argument))
	}
	return strings.Join(output, " ")
}

func builtinConsole_log(call FunctionCall) Value {mydebug.INFO()
	fmt.Fprintln(os.Stdout, formatForConsole(call.ArgumentList))
	return Value{}
}

func builtinConsole_error(call FunctionCall) Value {mydebug.INFO()
	fmt.Fprintln(os.Stdout, formatForConsole(call.ArgumentList))
	return Value{}
}

// Nothing happens.
func builtinConsole_dir(call FunctionCall) Value {mydebug.INFO()
	return Value{}
}

func builtinConsole_time(call FunctionCall) Value {mydebug.INFO()
	return Value{}
}

func builtinConsole_timeEnd(call FunctionCall) Value {mydebug.INFO()
	return Value{}
}

func builtinConsole_trace(call FunctionCall) Value {mydebug.INFO()
	return Value{}
}

func builtinConsole_assert(call FunctionCall) Value {mydebug.INFO()
	return Value{}
}

func (runtime *_runtime) newConsole() *_object {mydebug.INFO()

	return newConsoleObject(runtime)
}
