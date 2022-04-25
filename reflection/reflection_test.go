package reflection_test

import (
	"encoding/json"
	"testing"

	"github.com/enorith/supports/reflection"
)

var ts = reflection.InterfaceType[json.Marshaler]()

func TestInterfaceType(t *testing.T) {
	t.Log(ts)
}
