package utils

import (
	"testing"
	"time"
)

type User struct {
	Height int64
}

type PointerTimeSource struct {
	Name      string
	Age       int
	Is        bool
	Time      *time.Time
	TestMap   map[string]string
	TestSlice []string
	Arr       [10]float32
	User
	User1     User
	User2     *User
	MyChannel chan User
	GetUser   func()
	TypeName  interface{}
	Price     string
}

type TimeSource struct {
	Name      string
	Age       int
	Is        bool
	Time      time.Time
	TestMap   map[string]string
	TestSlice []string
	Arr       [10]float32
	User
	User1     User
	User2     *User
	MyChannel chan User
	GetUser   func()
	TypeName  interface{}
	Price     string
}

type StringTimeSource struct {
	Name      string
	Age       int
	Is        bool
	Time      string
	TestMap   map[string]string
	TestSlice []string
	Arr       [10]float32
	User
	User1     User
	User2     *User
	MyChannel chan User
	GetUser   func()
	TypeName  interface{}
	Price     string
}

var pointerTimeSource PointerTimeSource
var timeSource TimeSource
var stringTimeSource StringTimeSource

type StringDestination struct {
	Name      string
	Age       int
	Is        bool
	Time      string
	TestMap   map[string]string
	TestSlice []string
	Arr       [10]float32
	User
	User1     User
	User2     *User
	MyChannel chan User
	GetUser   func()
	TypeName  interface{}
	Price     string
}

type TimeDestination struct {
	Name      string
	Age       int
	Is        bool
	Time      time.Time
	TestMap   map[string]string
	TestSlice []string
	Arr       [10]float32
	User
	User1     User
	User2     *User
	MyChannel chan User
	GetUser   func()
	TypeName  interface{}
	Price     string
}

type PointerTimeDestination struct {
	Name      string
	Age       int
	Is        bool
	Time      *time.Time
	TestMap   map[string]string
	TestSlice []string
	Arr       [10]float32
	User
	User1     User
	User2     *User
	MyChannel chan User
	GetUser   func()
	TypeName  interface{}
	Price     string
}

func init() {
	//time.Local = time.FixedZone("UTC", 0)

	testMap := make(map[string]string)
	testMap["name"] = "berners"
	testSlice := make([]string, 1)
	testSlice = append(testSlice, "a")
	myChannel := make(chan User, 1)
	myChannel <- User{Height: 120}
	tt := time.Now()

	pointerTimeSource = PointerTimeSource{
		Name:      "John",
		Age:       30,
		Is:        true,
		Time:      &tt,
		TestMap:   testMap,
		User:      User{Height: 100},
		User1:     User{Height: 110},
		User2:     &User{Height: 120},
		MyChannel: myChannel,
		GetUser:   func() {},
		TypeName:  1,
	}

	timeSource = TimeSource{
		Name:      "John",
		Age:       30,
		Is:        true,
		Time:      tt,
		TestMap:   testMap,
		User:      User{Height: 100},
		User1:     User{Height: 110},
		User2:     &User{Height: 120},
		MyChannel: myChannel,
		GetUser:   func() {},
		TypeName:  1,
	}

	stringTimeSource = StringTimeSource{
		Name:      "John",
		Age:       30,
		Is:        true,
		Time:      "2023-03-23T14:39:14.8562572+08:00",
		TestMap:   testMap,
		User:      User{Height: 100},
		User1:     User{Height: 110},
		User2:     &User{Height: 120},
		MyChannel: myChannel,
		GetUser:   func() {},
		TypeName:  1,
	}
}

// go test -v mappingStruct_test.go mappingStruct.go

// *time.Time -> string
// *time.Time -> *time.Time
// *time.Time -> time.Time

// go test -v -run Test_mapStruct_pointerTimeToString mappingStruct_test.go mappingStruct.go
func Test_mapStruct_pointerTimeToString(t *testing.T) {
	t.Log(pointerTimeSource)

	var dest1 StringDestination
	MappingStruct(&pointerTimeSource, &dest1)
	t.Log(dest1)

	var dest2 StringDestination
	MappingStruct(pointerTimeSource, &dest2)
	t.Log(dest2)
}

// go test -v -run Test_mapStruct_pointerTimeToPointerTime mappingStruct_test.go mappingStruct.go
func Test_mapStruct_pointerTimeToPointerTime(t *testing.T) {
	t.Log(pointerTimeSource)

	var dest1 PointerTimeDestination
	MappingStruct(&pointerTimeSource, &dest1)
	t.Log(dest1)

	var dest2 PointerTimeDestination
	MappingStruct(pointerTimeSource, &dest2)
	t.Log(dest2)
}

// go test -v -run Test_mapStruct_pointerTimeToTime mappingStruct_test.go mappingStruct.go
func Test_mapStruct_pointerTimeToTime(t *testing.T) {
	t.Log(pointerTimeSource)

	var dest1 TimeDestination
	MappingStruct(&pointerTimeSource, &dest1)
	t.Log(dest1)

	var dest2 TimeDestination
	MappingStruct(pointerTimeSource, &dest2)
	t.Log(dest2)
}

// ----------------------------------------------------------------------------------------

// time.Time -> string
// time.Time -> *time.Time
// time.Time -> time.Time

// go test -v -run Test_mapStruct_timeToString mappingStruct_test.go mappingStruct.go
func Test_mapStruct_timeToString(t *testing.T) {
	t.Log(timeSource)

	var dest1 StringDestination
	MappingStruct(&timeSource, &dest1)
	t.Log(dest1)

	var dest2 StringDestination
	MappingStruct(timeSource, &dest2)
	t.Log(dest2)
}

// go test -v -run Test_mapStruct_timeToPointerTime mappingStruct_test.go mappingStruct.go
func Test_mapStruct_timeToPointerTime(t *testing.T) {
	t.Log(timeSource)

	var dest1 PointerTimeDestination
	MappingStruct(&timeSource, &dest1)
	t.Log(dest1)

	var dest2 PointerTimeDestination
	MappingStruct(timeSource, &dest2)
	t.Log(dest2)
}

// go test -v -run Test_mapStruct_timeToTime mappingStruct_test.go mappingStruct.go
func Test_mapStruct_timeToTime(t *testing.T) {
	t.Log(timeSource)

	var dest1 TimeDestination
	MappingStruct(&timeSource, &dest1)
	t.Log(dest1)

	var dest2 TimeDestination
	MappingStruct(timeSource, &dest2)
	t.Log(dest2)
}

// ----------------------------------------------------------------------------------------

// string -> string
// string -> *time.Time
// string -> time.Time

// go test -v -run Test_mapStruct_stringTimeToString mappingStruct_test.go mappingStruct.go
func Test_mapStruct_stringTimeToString(t *testing.T) {
	t.Log(stringTimeSource)

	var dest1 StringDestination
	MappingStruct(&stringTimeSource, &dest1)
	t.Log(dest1)

	var dest2 StringDestination
	MappingStruct(stringTimeSource, &dest2)
	t.Log(dest2)
}

// go test -v -run Test_mapStruct_stringTimeToPointerTime mappingStruct_test.go mappingStruct.go
func Test_mapStruct_stringTimeToPointerTime(t *testing.T) {
	t.Log(stringTimeSource)

	var dest1 PointerTimeDestination
	MappingStruct(&stringTimeSource, &dest1)
	t.Log(dest1)

	var dest2 PointerTimeDestination
	MappingStruct(stringTimeSource, &dest2)
	t.Log(dest2)
}

// go test -v -run Test_mapStruct_stringTimeToTime mappingStruct_test.go mappingStruct.go
func Test_mapStruct_stringTimeToTime(t *testing.T) {
	t.Log(stringTimeSource)

	var dest1 TimeDestination
	MappingStruct(&stringTimeSource, &dest1)
	t.Log(dest1)

	var dest2 TimeDestination
	MappingStruct(stringTimeSource, &dest2)
	t.Log(dest2)
}

// ----------------------------------------------------------------------------------------

type GoodsSource struct {
	Uint   string
	Uint8  string
	Uint16 string
	Uint32 string
	Uint64 string
	Int    string
	Int8   string
	Int16  string
	Int32  string
	Int64  string
}

type GoodsPointerSource struct {
	Uint   *string
	Uint8  *string
	Uint16 *string
	Uint32 *string
	Uint64 *string
	Int    *string
	Int8   *string
	Int16  *string
	Int32  *string
	Int64  *string
}

type GoodsDestination struct {
	Uint   uint
	Uint8  uint8
	Uint16 uint16
	Uint32 uint32
	Uint64 uint64
	Int    int
	Int8   int8
	Int16  int16
	Int32  int32
	Int64  int64
}

type GoodsPointerDestination struct {
	Uint   *uint
	Uint8  *uint8
	Uint16 *uint16
	Uint32 *uint32
	Uint64 *uint64
	Int    *int
	Int8   *int8
	Int16  *int16
	Int32  *int32
	Int64  *int64
}

// uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64
// go test -v -run Test_mapStruct_stringPointerToNumber mappingStruct_test.go mappingStruct.go
func Test_mapStruct_stringPointerToNumber(t *testing.T) {
	var goodsSource GoodsPointerSource
	Uint := "1"
	Uint8 := "2"
	Uint16 := "3"
	Uint32 := "4"
	Uint64 := "5"
	Int := "11"
	Int8 := "12"
	Int16 := "13"
	Int32 := "14"
	Int64 := "15"
	goodsSource.Uint = &Uint
	goodsSource.Uint8 = &Uint8
	goodsSource.Uint16 = &Uint16
	goodsSource.Uint32 = &Uint32
	goodsSource.Uint64 = &Uint64
	goodsSource.Int = &Int
	goodsSource.Int8 = &Int8
	goodsSource.Int16 = &Int16
	goodsSource.Int32 = &Int32
	goodsSource.Int64 = &Int64

	var goodsDestination1 GoodsDestination
	MappingStruct(goodsSource, &goodsDestination1)
	t.Log(goodsDestination1)

	var goodsDestination2 GoodsDestination
	MappingStruct(&goodsSource, &goodsDestination2)
	t.Log(goodsDestination2)
}

// go test -v -run Test_mapStruct_stringToNumber mappingStruct_test.go mappingStruct.go
func Test_mapStruct_stringToNumber(t *testing.T) {
	var goodsSource GoodsSource
	goodsSource.Uint = "1"
	goodsSource.Uint8 = "2"
	goodsSource.Uint16 = "3"
	goodsSource.Uint32 = "4"
	goodsSource.Uint64 = "5"
	goodsSource.Int = "11"
	goodsSource.Int8 = "12"
	goodsSource.Int16 = "13"
	goodsSource.Int32 = "14"
	goodsSource.Int64 = "15"

	var goodsDestination1 GoodsDestination
	MappingStruct(goodsSource, &goodsDestination1)
	t.Log(goodsDestination1)

	var goodsDestination2 GoodsDestination
	MappingStruct(&goodsSource, &goodsDestination2)
	t.Log(goodsDestination2)
}

// go test -v -run Test_mapStruct_numberToString mappingStruct_test.go mappingStruct.go
func Test_mapStruct_numberToString(t *testing.T) {
	var goodsDestination1 GoodsDestination
	goodsDestination1.Uint = 1
	goodsDestination1.Uint8 = 2
	goodsDestination1.Uint16 = 3
	goodsDestination1.Uint32 = 4
	goodsDestination1.Uint64 = 5
	goodsDestination1.Int = 11
	goodsDestination1.Int8 = 12
	goodsDestination1.Int16 = 13
	goodsDestination1.Int32 = 14
	goodsDestination1.Int64 = 15
	var goodsSource GoodsSource
	MappingStruct(&goodsDestination1, &goodsSource)
	t.Log(goodsSource)
}

// go test -v -run Test_mapStruct_numberPointerToString mappingStruct_test.go mappingStruct.go
func Test_mapStruct_numberPointerToString(t *testing.T) {
	var goodsPointerDestination GoodsPointerDestination
	var Uint uint = 1
	var Uint8 uint8 = 2
	var Uint16 uint16 = 3
	var Uint32 uint32 = 4
	var Uint64 uint64 = 5
	var Int = 11
	var Int8 int8 = 12
	var Int16 int16 = 13
	var Int32 int32 = 14
	var Int64 int64 = 15
	goodsPointerDestination.Uint = &Uint
	goodsPointerDestination.Uint8 = &Uint8
	goodsPointerDestination.Uint16 = &Uint16
	goodsPointerDestination.Uint32 = &Uint32
	goodsPointerDestination.Uint64 = &Uint64
	goodsPointerDestination.Int = &Int
	goodsPointerDestination.Int8 = &Int8
	goodsPointerDestination.Int16 = &Int16
	goodsPointerDestination.Int32 = &Int32
	goodsPointerDestination.Int64 = &Int64

	var goodsSource1 GoodsSource
	MappingStruct(&goodsPointerDestination, &goodsSource1)
	t.Log(goodsSource1)

	var goodsSource2 GoodsSource
	MappingStruct(goodsPointerDestination, &goodsSource2)
	t.Log(goodsSource2)
}
