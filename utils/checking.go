package utils

import (
	"reflect"
	"regexp"
	"strconv"

	"github.com/lucifinil-long/stores/proto"
	"github.com/mkideal/log"
)

// CheckDataValid checks whether data struct value obeys orm defines
// @param val is struct object to be checked
// @param fields is field names which want to be checked
// @return (error field name, error) if value does not obey orm defines; otherwise return ("", nil)
// @NOTE this will not check sub field if field is struct type
func CheckDataValid(val interface{}, fields ...string) (string, error) {
	vt := reflect.TypeOf(val)
	v := reflect.ValueOf(val)

	for vt.Kind() == reflect.Ptr {
		v = v.Elem()
		vt = v.Type()
	}
	if vt.Kind() != reflect.Struct {
		log.Fatal("utils.CheckDataValid get %v, source type: %T", proto.ErrNotStructType, val)
	}

	notNullRegexp, _ := regexp.Compile("not null")
	autoIncrRegexp, _ := regexp.Compile("autoincr")
	defaultRegexp, _ := regexp.Compile("default")
	varcharRegexp, _ := regexp.Compile("VARCHAR\\((\\d+)\\)")
	for i := 0; i < vt.NumField(); i++ {
		field := vt.Field(i)
		if len(fields) > 0 { // if set fields, only check matched fields
			found := false
			for _, f := range fields {
				if f == field.Name {
					found = true
					break
				}
			}

			if !found {
				continue
			}
		}

		tag := field.Tag.Get("xorm")
		if len(tag) == 0 || tag == "-" {
			continue
		}

		notNull := notNullRegexp.MatchString(tag)
		autoIncr := autoIncrRegexp.MatchString(tag)
		defaultVal := defaultRegexp.MatchString(tag)
		matcheds := varcharRegexp.FindSubmatch([]byte(tag))
		matchedArr := make([]string, 0, len(matcheds))
		for _, matched := range matcheds {
			matchedArr = append(matchedArr, string(matched))
		}

		// in this way, we have not null but without any autoincr or default flag, so we need to check value
		if notNull && !autoIncr && !defaultVal && v.Field(i).Kind() != reflect.Struct && v.Field(i).IsNil() {
			return field.Name, proto.ErrInvalidValueForNutNullField
		}

		if len(matchedArr) >= 2 {
			maxLen, _ := strconv.Atoi(matchedArr[1])

			if maxLen <= 0 {
				log.Fatal("utils.CheckDataValid get wrong varchar length, parsed values: %v", matchedArr)
			}

			str := v.Field(i).String()
			if len(str) > maxLen {
				return field.Name, proto.ErrDataTooLong
			}
		}
	}

	return "", nil
}
