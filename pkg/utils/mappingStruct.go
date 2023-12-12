package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func TimeStrToTime(timeStr string) (time.Time, error) {
	layouts := []string{
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
		time.DateTime,
		time.DateOnly,
		time.TimeOnly,
	}
	var t, _ = time.Parse("", "")
	var ts = t.String()

	if timeStr == "" {
		return t, errors.New("empty datetime are meaningless")
	}

	for _, layout := range layouts {
		tt, err := time.Parse(layout, timeStr)
		if err != nil {
			continue
		}
		t = tt
		break
	}
	// No matching time format found
	if t.String() == ts {
		return t, errors.New("datetime string format error")
	}
	return t, nil
}

func mappingStruct(l interface{}, r interface{}, opt MappingStructOption) error {
	getType, getValue := reflect.TypeOf(l), reflect.ValueOf(l)
	_, setValue := reflect.TypeOf(r), reflect.ValueOf(r)

	// right check Pointer
	if setValue.Kind() != reflect.Ptr {
		return errors.New("r not is pointer")
	}
	if setValue.IsNil() {
		return errors.New("r is empty pointer")
	}

	if getValue.Kind() == reflect.Ptr && !getValue.IsNil() {
		v := setValue.Elem()

		for i := 0; i < getType.Elem().NumField(); i++ {
			// 左字段
			leftField := getType.Elem().Field(i)
			// 左值
			leftValue := getValue.Elem().Field(i)

			// 左边的值空指针
			if leftValue.Kind() == reflect.Ptr && leftValue.IsNil() {
				//fmt.Println("左侧空指针：", field.Name, "---", getValue.Elem().Field(i))
				continue
			}

			// 左边的 字段 不在右边，就跳过
			rightField := v.FieldByName(leftField.Name)
			if rightField.String() == "<invalid Value>" {
				//fmt.Println("左边的字段不在右边：", field.Name, "---", getValue.Elem().Field(i))
				continue
			}

			// 左侧类型
			leftFieldType := leftField.Type.String()
			// 右侧类型
			rightFieldType := rightField.Type().String()

			var stringToNumber = func(str string) {
				if str == "" {
					str = "0"
				}
				num, err := strconv.ParseUint(str, 10, 64)
				if err != nil {
					return
				}

				if rightFieldType == "uint" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(uint(num)))
				}
				if rightFieldType == "*uint" {
					t := uint(num)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&t))
				}
				if rightFieldType == "uint8" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(uint8(num)))
				}
				if rightFieldType == "*uint8" {
					t := uint8(num)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&t))
				}
				if rightFieldType == "uint16" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(uint16(num)))
				}
				if rightFieldType == "*uint16" {
					t := uint16(num)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&t))
				}
				if rightFieldType == "uint32" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(uint32(num)))
				}
				if rightFieldType == "*uint32" {
					t := uint32(num)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(t))
				}
				if rightFieldType == "uint64" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(num))
				}
				if rightFieldType == "*uint64" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&num))
				}
				if rightFieldType == "int" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(int(num)))
				}
				if rightFieldType == "*int" {
					t := int(num)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&t))
				}
				if rightFieldType == "int8" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(int8(num)))
				}
				if rightFieldType == "*int8" {
					t := int8(num)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&t))
				}
				if rightFieldType == "int16" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(int16(num)))
				}
				if rightFieldType == "*int16" {
					t := int16(num)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&t))
				}
				if rightFieldType == "int32" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(int32(num)))
				}
				if rightFieldType == "*int32" {
					t := int32(num)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&t))
				}
				if rightFieldType == "int64" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(int64(num)))
				}
				if rightFieldType == "*int64" {
					t := int64(num)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(t))
				}
			}

			switch leftFieldType {
			case "*time.Time":
				// *time.Time -> *time.Time
				// *time.Time -> time.Time
				// *time.Time -> string
				// *time.Time -> *string
				if rightFieldType == "*time.Time" {
					v.FieldByName(leftField.Name).Set(leftValue)
				}

				_time := leftValue.Interface().(*time.Time)
				if _time == nil {
					continue
				}
				if rightFieldType == "time.Time" {
					// *time.Time -> string -> time.Time
					t, err := TimeStrToTime(_time.Format(opt.FormatLayout))
					if err != nil {
						continue
					}
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(t))
				}

				if rightFieldType == "string" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(_time.Format(opt.FormatLayout)))
				}
				if rightFieldType == "*string" {
					t := _time.Format(opt.FormatLayout)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&t))
				}

			case "time.Time":
				// time.Time -> time.Time
				// time.Time -> *time.Time
				// time.Time -> string
				// time.Time -> *string
				if rightFieldType == "time.Time" {
					v.FieldByName(leftField.Name).Set(leftValue)
				}
				_time := leftValue.Interface().(time.Time)
				if rightFieldType == "*time.Time" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&_time))
				}
				if rightFieldType == "string" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(_time.Format(opt.FormatLayout)))
				}
				if rightFieldType == "*string" {
					t := _time.Format(opt.FormatLayout)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&t))
				}
			case "*string":
				// *string -> *string
				// *string -> string
				// *string -> *time.Time
				// *string -> time.Time

				// *string -> uint| uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64
				// *string -> *uint| *uint8 | *uint16 | *uint32 | *uint64 | *int | *int8 | *int16 | *int32 | *int64
				if rightFieldType == "*string" {
					v.FieldByName(leftField.Name).Set(leftValue)
				}

				str := leftValue.Interface().(*string)
				if str == nil {
					continue
				}
				if rightFieldType == "string" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(*str))
				}
				if rightFieldType == "*time.Time" {
					_time, err := TimeStrToTime(*str)
					if err != nil {
						continue
					}
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&_time))
				}
				if rightFieldType == "time.Time" {
					_time, err := TimeStrToTime(*str)
					if err != nil {
						continue
					}
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(_time))
				}

				stringToNumber(*str)
			case "string":
				// string -> string
				// string -> *string
				// string -> *time.Time
				// string -> time.Time

				// string -> uint| uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64
				// string -> *uint| *uint8 | *uint16 | *uint32 | *uint64 | *int | *int8 | *int16 | *int32 | *int64
				if rightFieldType == "string" {
					v.FieldByName(leftField.Name).Set(leftValue)
				}
				if rightFieldType == "*string" {
					str := leftValue.String()
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&str))
				}
				if rightFieldType == "*time.Time" {
					_time, err := TimeStrToTime(leftValue.String())
					if err != nil {
						continue
					}
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&_time))
				}
				if rightFieldType == "time.Time" {
					_time, err := TimeStrToTime(leftValue.String())
					if err != nil {
						continue
					}
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(_time))
				}

				stringToNumber(leftValue.String())
			case "uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64":
				// number -> number
				// number -> string
				// number -> *string
				if leftFieldType == rightFieldType {
					v.FieldByName(leftField.Name).Set(leftValue)
				}

				var numberToString = func(str string) {
					if str == "0" {
						str = ""
					}
					if rightFieldType == "string" {
						v.FieldByName(leftField.Name).Set(reflect.ValueOf(str))
					}
					if rightFieldType == "*string" {
						v.FieldByName(leftField.Name).Set(reflect.ValueOf(&str))
					}
				}
				str := fmt.Sprintf("%d", leftValue.Interface())
				numberToString(str)
			case "*uint", "*uint8", "*uint16", "*uint32", "*uint64", "*int", "*int8", "*int16", "*int32", "*int64":
				// *number -> *number
				// *number -> string
				// *number -> *string
				if leftFieldType == rightFieldType {
					v.FieldByName(leftField.Name).Set(leftValue)
				}

				var numberToString = func(str string) {
					if str == "0" {
						str = ""
					}
					if rightFieldType == "string" {
						v.FieldByName(leftField.Name).Set(reflect.ValueOf(str))
					}
					if rightFieldType == "*string" {
						v.FieldByName(leftField.Name).Set(reflect.ValueOf(&str))
					}
				}

				if leftFieldType == "*uint" {
					num := leftValue.Interface().(*uint)
					if num == nil {
						continue
					}
					str := fmt.Sprintf("%d", *num)
					numberToString(str)
				}
				if leftFieldType == "*uint8" {
					num := leftValue.Interface().(*uint8)
					if num == nil {
						continue
					}
					str := fmt.Sprintf("%d", *num)
					numberToString(str)
				}
				if leftFieldType == "*uint16" {
					num := leftValue.Interface().(*uint16)
					if num == nil {
						continue
					}
					str := fmt.Sprintf("%d", *num)
					numberToString(str)
				}
				if leftFieldType == "*uint32" {
					num := leftValue.Interface().(*uint32)
					if num == nil {
						continue
					}
					str := fmt.Sprintf("%d", *num)
					numberToString(str)
				}
				if leftFieldType == "*uint64" {
					num := leftValue.Interface().(*uint64)
					if num == nil {
						continue
					}
					str := fmt.Sprintf("%d", *num)
					numberToString(str)
				}
				if leftFieldType == "*int" {
					num := leftValue.Interface().(*int)
					if num == nil {
						continue
					}
					str := fmt.Sprintf("%d", *num)
					numberToString(str)
				}
				if leftFieldType == "*int8" {
					num := leftValue.Interface().(*int8)
					if num == nil {
						continue
					}
					str := fmt.Sprintf("%d", *num)
					numberToString(str)
				}
				if leftFieldType == "*int16" {
					num := leftValue.Interface().(*int16)
					if num == nil {
						continue
					}
					str := fmt.Sprintf("%d", *num)
					numberToString(str)
				}
				if leftFieldType == "*int32" {
					num := leftValue.Interface().(*int32)
					if num == nil {
						continue
					}
					str := fmt.Sprintf("%d", *num)
					numberToString(str)
				}
				if leftFieldType == "*int64" {
					num := leftValue.Interface().(*int64)
					if num == nil {
						continue
					}
					str := fmt.Sprintf("%d", *num)
					numberToString(str)
				}
			default:
				// 两侧数据类型一至才可以映射
				if leftFieldType == rightFieldType &&
					leftValue.Kind().String() != "ptr" &&
					leftValue.Kind().String() != "array" &&
					leftValue.Kind().String() != "struct" &&
					leftValue.Kind().String() != "chan" &&
					leftValue.Kind().String() != "func" &&
					leftValue.Kind().String() != "slice" &&
					leftValue.Kind().String() != "interface" &&
					leftValue.Kind().String() != "map" {
					v.FieldByName(leftField.Name).Set(leftValue)
				}
			}
		}
	}

	if getValue.Kind() == reflect.Struct {
		v := setValue.Elem()
		for i := 0; i < getType.NumField(); i++ {
			// 左字段
			leftField := getType.Field(i)
			// 左值
			leftValue := getValue.Field(i)

			// 左边的值空指针
			if leftValue.Kind() == reflect.Ptr && leftValue.IsNil() {
				//fmt.Println("左侧空指针：", field.Name, "---", getValue.Field(i))
				continue
			}
			// 左边的 字段 不在右边，就跳过
			if v.FieldByName(leftField.Name).String() == "<invalid Value>" {
				//fmt.Println("左边的字段不在右边：", field.Name, "---", getValue.Field(i))
				continue
			}

			// 左侧类型
			leftFieldType := leftField.Type.String()
			// 右侧类型
			rightFieldType := v.FieldByName(leftField.Name).Type().String()

			var stringToNumber = func(str string) {
				if str == "" {
					str = "0"
				}
				num, err := strconv.ParseUint(str, 10, 64)
				if err != nil {
					return
				}

				if rightFieldType == "uint" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(uint(num)))
				}
				if rightFieldType == "*uint" {
					t := uint(num)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&t))
				}
				if rightFieldType == "uint8" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(uint8(num)))
				}
				if rightFieldType == "*uint8" {
					t := uint8(num)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&t))
				}
				if rightFieldType == "uint16" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(uint16(num)))
				}
				if rightFieldType == "*uint16" {
					t := uint16(num)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&t))
				}
				if rightFieldType == "uint32" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(uint32(num)))
				}
				if rightFieldType == "*uint32" {
					t := uint32(num)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(t))
				}
				if rightFieldType == "uint64" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(num))
				}
				if rightFieldType == "*uint64" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&num))
				}
				if rightFieldType == "int" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(int(num)))
				}
				if rightFieldType == "*int" {
					t := int(num)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&t))
				}
				if rightFieldType == "int8" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(int8(num)))
				}
				if rightFieldType == "*int8" {
					t := int8(num)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&t))
				}
				if rightFieldType == "int16" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(int16(num)))
				}
				if rightFieldType == "*int16" {
					t := int16(num)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&t))
				}
				if rightFieldType == "int32" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(int32(num)))
				}
				if rightFieldType == "*int32" {
					t := int32(num)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&t))
				}
				if rightFieldType == "int64" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(int64(num)))
				}
				if rightFieldType == "*int64" {
					t := int64(num)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(t))
				}
			}

			switch leftFieldType {
			// 左侧是 *time.Time 类型
			case "*time.Time":
				// *time.Time -> *time.Time
				// *time.Time -> time.Time
				// *time.Time -> *string
				// *time.Time -> string
				if rightFieldType == "*time.Time" {
					v.FieldByName(leftField.Name).Set(leftValue)
				}
				_time := leftValue.Interface().(*time.Time)
				if _time == nil {
					continue
				}
				if rightFieldType == "time.Time" {
					// *time.Time -> string -> time.Time
					t, err := TimeStrToTime(_time.Format(opt.FormatLayout))
					if err != nil {
						continue
					}
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(t))
				}
				if rightFieldType == "*string" {
					str := _time.Format(opt.FormatLayout)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&str))
				}
				if rightFieldType == "string" {
					str := _time.Format(opt.FormatLayout)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(str))
				}
			case "time.Time":
				// time.Time -> time.Time
				// time.Time -> *time.Time
				// time.Time -> string
				// time.Time -> *string
				if rightFieldType == "time.Time" {
					v.FieldByName(leftField.Name).Set(leftValue)
				}
				_time := leftValue.Interface().(time.Time)
				if rightFieldType == "*time.Time" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&_time))
				}
				if rightFieldType == "string" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(_time.Format(opt.FormatLayout)))
				}
				if rightFieldType == "*string" {
					t := _time.Format(opt.FormatLayout)
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&t))
				}
			case "*string":
				// *string -> *string
				// *string -> string
				// *string -> *time.Time
				// *string -> time.Time

				// *string -> uint| uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64
				// *string -> *uint| *uint8 | *uint16 | *uint32 | *uint64 | *int | *int8 | *int16 | *int32 | *int64
				if rightFieldType == "*string" {
					v.FieldByName(leftField.Name).Set(leftValue)
				}

				str := leftValue.Interface().(*string)
				if str == nil {
					continue
				}
				if rightFieldType == "string" {
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(*str))
				}
				if rightFieldType == "*time.Time" {
					_time, err := TimeStrToTime(*str)
					if err != nil {
						continue
					}
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&_time))
				}
				if rightFieldType == "time.Time" {
					_time, err := TimeStrToTime(*str)
					if err != nil {
						continue
					}
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(_time))
				}

				stringToNumber(*str)
			case "string":
				// string -> string
				// string -> *string
				// string -> *time.Time
				// string -> time.Time

				// string -> uint| uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64
				// string -> *uint| *uint8 | *uint16 | *uint32 | *uint64 | *int | *int8 | *int16 | *int32 | *int64
				if rightFieldType == "string" {
					v.FieldByName(leftField.Name).Set(leftValue)
				}
				if rightFieldType == "*string" {
					str := leftValue.String()
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&str))
				}
				if rightFieldType == "*time.Time" {
					_time, err := TimeStrToTime(leftValue.String())
					if err != nil {
						continue
					}
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(&_time))
				}
				if rightFieldType == "time.Time" {
					_time, err := TimeStrToTime(leftValue.String())
					if err != nil {
						continue
					}
					v.FieldByName(leftField.Name).Set(reflect.ValueOf(_time))
				}

				stringToNumber(leftValue.String())
			case "uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64":
				// number -> number
				// number -> string
				// number -> *string
				if leftFieldType == rightFieldType {
					v.FieldByName(leftField.Name).Set(leftValue)
				}

				var numberToString = func(str string) {
					if str == "0" {
						str = ""
					}
					if rightFieldType == "string" {
						v.FieldByName(leftField.Name).Set(reflect.ValueOf(str))
					}
					if rightFieldType == "*string" {
						v.FieldByName(leftField.Name).Set(reflect.ValueOf(&str))
					}
				}
				str := fmt.Sprintf("%d", leftValue.Interface())
				numberToString(str)
			case "*uint", "*uint8", "*uint16", "*uint32", "*uint64", "*int", "*int8", "*int16", "*int32", "*int64":
				// *number -> *number
				// *number -> string
				// *number -> *string
				if leftFieldType == rightFieldType {
					v.FieldByName(leftField.Name).Set(leftValue)
				}

				var numberToString = func(str string) {
					if str == "0" {
						str = ""
					}
					if rightFieldType == "string" {
						v.FieldByName(leftField.Name).Set(reflect.ValueOf(str))
					}
					if rightFieldType == "*string" {
						v.FieldByName(leftField.Name).Set(reflect.ValueOf(&str))
					}
				}

				if leftFieldType == "*uint" {
					num := leftValue.Interface().(*uint)
					if num == nil {
						continue
					}
					str := fmt.Sprintf("%d", *num)
					numberToString(str)
				}
				if leftFieldType == "*uint8" {
					num := leftValue.Interface().(*uint8)
					if num == nil {
						continue
					}
					str := fmt.Sprintf("%d", *num)
					numberToString(str)
				}
				if leftFieldType == "*uint16" {
					num := leftValue.Interface().(*uint16)
					if num == nil {
						continue
					}
					str := fmt.Sprintf("%d", *num)
					numberToString(str)
				}
				if leftFieldType == "*uint32" {
					num := leftValue.Interface().(*uint32)
					if num == nil {
						continue
					}
					str := fmt.Sprintf("%d", *num)
					numberToString(str)
				}
				if leftFieldType == "*uint64" {
					num := leftValue.Interface().(*uint64)
					if num == nil {
						continue
					}
					str := fmt.Sprintf("%d", *num)
					numberToString(str)
				}
				if leftFieldType == "*int" {
					num := leftValue.Interface().(*int)
					if num == nil {
						continue
					}
					str := fmt.Sprintf("%d", *num)
					numberToString(str)
				}
				if leftFieldType == "*int8" {
					num := leftValue.Interface().(*int8)
					if num == nil {
						continue
					}
					str := fmt.Sprintf("%d", *num)
					numberToString(str)
				}
				if leftFieldType == "*int16" {
					num := leftValue.Interface().(*int16)
					if num == nil {
						continue
					}
					str := fmt.Sprintf("%d", *num)
					numberToString(str)
				}
				if leftFieldType == "*int32" {
					num := leftValue.Interface().(*int32)
					if num == nil {
						continue
					}
					str := fmt.Sprintf("%d", *num)
					numberToString(str)
				}
				if leftFieldType == "*int64" {
					num := leftValue.Interface().(*int64)
					if num == nil {
						continue
					}
					str := fmt.Sprintf("%d", *num)
					numberToString(str)
				}
			default:
				// 两侧数据类型一至才可以映射
				if leftFieldType == rightFieldType &&
					leftValue.Kind().String() != "ptr" &&
					leftValue.Kind().String() != "array" &&
					leftValue.Kind().String() != "struct" &&
					leftValue.Kind().String() != "chan" &&
					leftValue.Kind().String() != "func" &&
					leftValue.Kind().String() != "slice" &&
					leftValue.Kind().String() != "interface" &&
					leftValue.Kind().String() != "map" {
					v.FieldByName(leftField.Name).Set(leftValue)
				}
			}
		}
	}
	return nil
}

type MappingStructOption struct {
	FormatLayout string
}

func MappingStruct(l interface{}, r interface{}) {
	err := mappingStruct(l, r, MappingStructOption{FormatLayout: time.RFC3339Nano})
	if err != nil {
		return
	}
}

func MappingStructMappingStructOption(l interface{}, r interface{}, opt MappingStructOption) {
	err := mappingStruct(l, r, opt)
	if err != nil {
		return
	}
}
