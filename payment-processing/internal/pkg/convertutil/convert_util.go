package convertutil

import "strconv"

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
