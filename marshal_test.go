package urlquery

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func testMarshalToValues(t *testing.T) {

	Convey("marshalToValues misc error - 测试 marshalToValues 基础类型", func() {
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

	Convey("readTag - 测试 readTag", func() {
		type pet struct {
			OwnerID string `url:"owner_id,omitempty" json:"ownerID"`
			Name    string `url:",omitempty"`
			Sex     int
		}

		p := pet{
			OwnerID: "Tencent",
			Name:    "Penguin",
			Sex:     1,
		}

		typ := reflect.TypeOf(p)

		ft := typ.Field(0)
		tg := readTag(&ft, "url")
		So(tg.Name(), ShouldEqual, "owner_id")
		So(tg.Has("omitempty"), ShouldBeTrue)

		ft = typ.Field(0)
		tg = readTag(&ft, "json")
		So(tg.Name(), ShouldEqual, "ownerID")
		So(tg.Has("omitempty"), ShouldBeFalse)

		ft = typ.Field(1)
		tg = readTag(&ft, "url")
		So(tg.Name(), ShouldEqual, "Name")
		So(tg.Has("omitempty"), ShouldBeTrue)

		ft = typ.Field(2)
		tg = readTag(&ft, "url")
		So(tg.Name(), ShouldEqual, "Sex")
		So(tg.Has("omitempty"), ShouldBeFalse)
	})

}

func testMarshal(t *testing.T) {
	type Pet struct {
		OwnerID string `url:"owner_id,omitempty" json:"ownerID"`
		Name    string `url:",omitempty"`
		Sex     int
	}

	p := Pet{
		OwnerID: "tencent",
		Name:    "Penguin",
		Sex:     1,
	}

	s, err := Marshal(&p)
	So(err, ShouldBeNil)

	t.Log(string(s))

	u, err := url.ParseQuery(string(s))
	So(err, ShouldBeNil)
	So(len(u), ShouldEqual, 3)
	So(u.Get("owner_id"), ShouldEqual, p.OwnerID)
	So(u.Get("Name"), ShouldEqual, p.Name)
	So(u.Get("Sex"), ShouldEqual, fmt.Sprintf("%d", p.Sex))
}
