package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/transform"
)

// JoinArray joins array members as a string
// @param sep is seperator between each member
// @param array is variable list
// @return 连接后的字符串
//
func JoinArray(sep string, array ...interface{}) string {
	strVals := make([]string, 0, len(array))
	for _, val := range array {
		strVals = append(strVals, fmt.Sprint(val))
	}

	return strings.Join(strVals, sep)
}

// Strings2JSON formats json string
func Strings2JSON(str string) string {
	rs := []rune(str)
	jsons := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			jsons += string(r)
		} else {
			jsons += "\\u" + strconv.FormatInt(int64(rint), 16) // json
		}
	}

	return jsons
}

// String2MD5 computes md5 string of specified string
func String2MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	rs := hex.EncodeToString(h.Sum(nil))

	return rs
}

// ConvertToUTF8 convert string from charset encode to utf8 charset
func ConvertToUTF8(source, charset string) (string, error) {
	switch {
	case strings.EqualFold("utf-8", charset):
		return source, nil
	case strings.EqualFold("iso-8859-1", charset):
		return source, nil
	case strings.EqualFold("us-ascii", charset):
		return source, nil
	default:
		enc, err := htmlindex.Get(charset)
		if err != nil {
			return "", err
		}
		in := bytes.NewReader([]byte(source))
		out := transform.NewReader(in, enc.NewDecoder())
		result, err := ioutil.ReadAll(out)
		if err != nil {
			return "", err
		}
		return string(result), nil
	}
}

// DecodeString decode string for specified encode format, current only support base64
func DecodeString(encode, content string) (string, error) {
	if strings.EqualFold("base64", encode) {
		decode, err := base64.StdEncoding.DecodeString(content)
		return string(decode), err
	}

	return content, nil
}
