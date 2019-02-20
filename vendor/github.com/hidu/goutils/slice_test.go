package utils
import "github.com/equalll/mydebug"

import (
	"fmt"
	"testing"
)

func TestIsInArray(t *testing.T) {mydebug.INFO()
	arr := []string{"a", "b", "c"}
	a(arr)
}

func a(b interface{}) {mydebug.INFO()
	fmt.Println(b)
}
