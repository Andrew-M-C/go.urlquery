package urlquery

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func test(t *testing.T, scene string, f func(*testing.T)) {
	if t.Failed() {
		return
	}
	Convey(scene, t, func() {
		f(t)
	})
}

func TestSlice(t *testing.T) {
	test(t, "marshalToValues", testMarshalToValues)
	test(t, "Marshal", testMarshal)
}
