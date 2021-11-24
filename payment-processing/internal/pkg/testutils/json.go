package testutils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func EqualJSON(actual interface{}, expected ...interface{}) string {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal(actual.([]byte), &o1)
	if err != nil {
		return fmt.Sprintf("Error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal(expected[0].([]byte), &o2)
	if err != nil {
		return fmt.Sprintf("Error mashalling string 2 :: %s", err.Error())
	}

	if !reflect.DeepEqual(o1, o2) {
		return fmt.Sprintf("actual \n %v \n and expectation \n %v \n not equal", o1, o2)
	}

	return ""
}

func KeysExists(actual interface{}, expected ...interface{}) string {
	var o1 map[string]interface{}

	var err error
	err = json.Unmarshal(actual.([]byte), &o1)
	if err != nil {
		return fmt.Sprintf("Error mashalling string 1 :: %s", err.Error())
	}

	for _, v := range expected {
		key := v.(string)
		if v, ok := o1[key]; !ok && v != nil {
			return fmt.Sprintf("Error key %s dont exists", key)
		}
	}

	return ""
}

func ToJSON(t *testing.T, data interface{}) []byte {
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
