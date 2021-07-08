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

func testMarshalStructInStruct(t *testing.T) {
	type sub struct {
		String string `url:"string"`
		Ints   []int  `url:"ints"`
	}
	type thing struct {
		Sub    sub  `url:"sub"`
		SubPtr *sub `url:"sub_ptr"`
		NilPtr *sub `url:"nil_ptr"`
	}

	th := &thing{
		Sub: sub{
			String: "string in sub",
			Ints:   []int{1, 3, 5, 7, 9},
		},
		SubPtr: &sub{
			String: "string in sub ptr",
			Ints:   []int{2, 4, 6, 8},
		},
		NilPtr: nil,
	}

	kv, err := marshalToValues(th)
	So(err, ShouldBeNil)

	t.Logf("result: %s", kv.Encode())

	// sub.xxx
	So(kv.Get("sub.string"), ShouldEqual, th.Sub.String)

	subInts := kv["sub.ints"]
	So(len(subInts), ShouldEqual, len(th.Sub.Ints))
	for i, s := range subInts {
		So(s, ShouldEqual, fmt.Sprint(th.Sub.Ints[i]))
	}

	// sub_ptr.xxx
	So(kv.Get("sub_ptr.string"), ShouldEqual, th.SubPtr.String)

	subPtrInts := kv["sub_ptr.ints"]
	So(len(subPtrInts), ShouldEqual, len(th.SubPtr.Ints))
	for i, s := range subPtrInts {
		So(s, ShouldEqual, fmt.Sprint(th.SubPtr.Ints[i]))
	}

	// nil_ptr.xxx
	_, exist := kv["nil_ptr"]
	So(exist, ShouldBeFalse)
}

func testMarshalMapInStruct(t *testing.T) {
	type thing struct {
		Intf   map[string]interface{} `url:"intf"`
		String map[string]string      `url:"string"`
	}

	th := &thing{
		Intf: map[string]interface{}{
			"string": "a string",
			"int":    12345,
			"bool":   true,
			"ints":   []int{2, 4, 6, 8, 10},
		},
		String: map[string]string{
			"A": "alpha",
			"B": "bravo",
			"C": "charlie",
		},
	}

	kv, err := marshalToValues(th)
	So(err, ShouldBeNil)

	t.Logf("result: %s", kv.Encode())

	So(kv.Get("intf.string"), ShouldEqual, th.Intf["string"])
	So(kv.Get("intf.int"), ShouldEqual, fmt.Sprint(th.Intf["int"]))
	So(kv.Get("intf.bool"), ShouldEqual, fmt.Sprint(th.Intf["bool"]))

	So(kv.Get("string.A"), ShouldEqual, th.String["A"])
	So(kv.Get("string.B"), ShouldEqual, th.String["B"])
	So(kv.Get("string.C"), ShouldEqual, th.String["C"])
}
