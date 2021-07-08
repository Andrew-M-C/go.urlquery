package urlquery

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
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
	type Dog struct {
		Pet
		Color string `url:"color,omitempty"`
	}

	p := Dog{
		Pet: Pet{
			OwnerID: "tencent",
			Name:    "Penguin",
			Sex:     1,
		},
		Color: "black",
	}

	s, err := Marshal(&p)
	So(err, ShouldBeNil)

	t.Log(string(s))

	u, err := url.ParseQuery(string(s))
	So(err, ShouldBeNil)
	So(len(u), ShouldEqual, 4)
	So(u.Get("owner_id"), ShouldEqual, p.OwnerID)
	So(u.Get("Name"), ShouldEqual, p.Name)
	So(u.Get("Sex"), ShouldEqual, fmt.Sprintf("%d", p.Sex))
	So(u.Get("color"), ShouldEqual, p.Color)
}

func testMarshalSlice(t *testing.T) {
	type thing struct {
		Strings  []string  `url:"strings"`
		Ints     []int     `url:"ints"`
		Uints    []uint    `url:"uints"`
		Bools    []bool    `url:"bools"`
		Floats   []float32 `url:"floats"`
		Empty    []int     `url:"empty"`
		HexArray [16]int8  `url:"hex_array"`
	}

	th := &thing{
		Strings:  []string{"A", "B", "C"},
		Ints:     []int{-1, -3, -5},
		Uints:    []uint{2, 4, 6},
		Bools:    []bool{true, false, true, false},
		Floats:   []float32{1.1, -2.2, 3.3},
		HexArray: [16]int8{1, 2, 3, 4, 5, 6, 7, 8, 9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF},
	}

	kv, err := marshalToValues(th)
	So(err, ShouldBeNil)

	t.Logf("result: %s", kv.Encode())

	strLst := kv["strings"]
	So(len(strLst), ShouldEqual, len(th.Strings))
	for i, s := range strLst {
		So(s, ShouldEqual, th.Strings[i])
	}

	intLst := kv["ints"]
	So(len(intLst), ShouldEqual, len(th.Ints))
	for i, s := range intLst {
		n, err := strconv.ParseInt(s, 10, 64)
		So(err, ShouldBeNil)
		So(n, ShouldEqual, th.Ints[i])
	}

	uintLst := kv["uints"]
	So(len(uintLst), ShouldEqual, len(th.Uints))
	for i, s := range uintLst {
		n, err := strconv.ParseUint(s, 10, 64)
		So(err, ShouldBeNil)
		So(n, ShouldEqual, th.Uints[i])
	}

	boolLst := kv["bools"]
	So(len(boolLst), ShouldEqual, len(th.Bools))
	for i, s := range boolLst {
		So(s, ShouldEqual, fmt.Sprint(th.Bools[i]))
	}

	floatLst := kv["floats"]
	So(len(floatLst), ShouldEqual, len(th.Floats))
	for i, s := range floatLst {
		f, err := strconv.ParseFloat(s, 32)
		So(err, ShouldBeNil)
		So(f, ShouldEqual, th.Floats[i])
	}

	hexArray := kv["hex_array"]
	So(len(hexArray), ShouldEqual, len(th.HexArray))
	for i, s := range hexArray {
		n, err := strconv.ParseInt(s, 10, 64)
		So(err, ShouldBeNil)
		So(n, ShouldEqual, th.HexArray[i])
	}

	_, exist := kv["empty"]
	So(exist, ShouldBeFalse)

}
