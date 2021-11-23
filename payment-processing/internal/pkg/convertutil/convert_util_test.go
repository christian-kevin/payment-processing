package convertutil

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	log "spenmo/payment-processing/payment-processing/pkg/logger"
	"testing"
)

func TestInflatePrice(t *testing.T) {
	log.InitializeLogger("TestInflatePrice")
	_ = gomock.NewController(t)

	tt := []struct {
		caseName     string
		price          string
		expectations func(t *testing.T, res int64)
	}{
		{
			caseName: "When price is zero",
			price:      "0",
			expectations: func(t *testing.T, res int64) {
				assert.Equal(t, int64(0), res)
			},
		},
		{
			caseName: "When price is integer",
			price:      "10000",
			expectations: func(t *testing.T, res int64) {
				assert.Equal(t, int64(100000000), res)
			},
		},
		{
			caseName: "When price is float",
			price:      "10.5",
			expectations: func(t *testing.T, res int64) {
				assert.Equal(t, int64(105000), res)
			},
		},
	}

	for _, l := range tt {
		t.Log(l.caseName)
		res, _ := InflatePrice(l.price)
		l.expectations(t, *res)
	}
}

func TestDeflatePrice(t *testing.T) {
	log.InitializeLogger("TestDeflatePrice")
	_ = gomock.NewController(t)

	tt := []struct {
		caseName     string
		price        int64
		expectations func(t *testing.T, res string)
	}{
		{
			caseName: "When price is zero",
			price:      0,
			expectations: func(t *testing.T, res string) {
				assert.Equal(t, "0", res)
			},
		},
		{
			caseName: "When price is integer",
			price:      100000000,
			expectations: func(t *testing.T, res string) {
				assert.Equal(t, "10000", res)
			},
		},
		{
			caseName: "When price is float",
			price:      105000,
			expectations: func(t *testing.T, res string) {
				assert.Equal(t, "10.5", res)
			},
		},
	}

	for _, l := range tt {
		t.Log(l.caseName)
		res := DeflatePrice(&l.price)
		l.expectations(t, res)
	}
}
