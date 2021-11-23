package cardutil

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateCardNumber() string {
	// ex: IIN for mastercard 51****
	first := 510000
	min := 100000000
	max := 999999999
	accNumber := rand.Intn(max - min) + min

	number1 := strconv.Itoa(first)
	number2 := strconv.Itoa(accNumber)

	digits := make([]int, 0, 0)

	for i := 0; i < len(number1) ; i++ {
		v, _ := strconv.Atoi(string(number1[i]))
		digits = append(digits, v)
	}

	for i := 0; i < len(number2) ; i++ {
		v, _ := strconv.Atoi(string(number2[i]))
		digits = append(digits, v)
	}

	doubleOddDigit := make([]int, 0, 0)
	for i := 0;i<len(digits);i++ {
		if (i+1)%2 != 0 {
			num := digits[i]*2
			if num > 9 {
				num = num - 9
			}

			doubleOddDigit = append(doubleOddDigit, num)
		} else {
			doubleOddDigit = append(doubleOddDigit, digits[i])
		}
	}

	checkSum := 0
	for i := 0;i<len(doubleOddDigit);i++ {
		checkSum += doubleOddDigit[i]
	}
	mod := checkSum%10
	if mod == 0 {
		checkSum = 0
	} else {
		checkSum = 10 - mod
	}

	res := number1+number2 + strconv.Itoa(checkSum)
	return res
}

func GenerateExpiryDate(country string) string {
	now := time.Now()
	after := now.AddDate(0, 24, 0) // add 2 years

	// For simplicity country is ID
	location, _ := time.LoadLocation("Asia/Jakarta")
	after = after.In(location)
	return after.Format("01/06")
}
