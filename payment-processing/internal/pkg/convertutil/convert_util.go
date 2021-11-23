package convertutil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func InflatePrice(price string) (*int64, error) {
	pF, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return nil, err
	}
	res := int64(pF * 10000)
	return &res, nil
}

func DeflatePrice(price *int64) string {
	pF := float64(*price) / 10000
	return strconv.FormatFloat(pF, 'f', -1, 64)
}

func MapRequestToModel(r *http.Request, model interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err := r.Body.Close(); err != nil {
		return err
	}
	if err := json.Unmarshal(body, model); err != nil {
		return err
	}
	return nil
}
