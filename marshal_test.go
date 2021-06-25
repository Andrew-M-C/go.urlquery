package urlquery

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func testMarshalToValues(t *testing.T) {
	Convey("marshalToValues misc error", func() {
		type data struct{}

		var err error
		var ptr *data

		_, err = marshalToValues(ptr)
		So(err, ShouldBeError)

		_, err = marshalToValues(123)
		So(err, ShouldBeError)
		t.Log(err)

		_, err = marshalToValues(data{})
		So(err, ShouldBeNil)

		_, err = marshalToValues(&data{})
		So(err, ShouldBeNil)
	})

}
