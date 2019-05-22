package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidators(t *testing.T) {
	t.Run("Integer Validator", func(t *testing.T) {
		assert := assert.New(t)

		var tests = []struct {
			input    int64
			kv       []string
			expected bool
		}{
			{122, []string{"max:123"}, true},
			{123, []string{"max:123"}, true},
			{124, []string{"max:123"}, false},
			{122, []string{"min:123"}, false},
			{123, []string{"min:123"}, true},
			{124, []string{"min:123"}, true},
			{84, []string{"max:123", "min:43"}, true},
			{42, []string{"max:123", "min:43"}, false},
			{124, []string{"max:123", "min:43"}, false},
		}

		for _, test := range tests {
			assert.Equal(test.expected, integerValidator(test.input, test.kv))
		}
	})

	t.Run("Float Validator", func(t *testing.T) {
		assert := assert.New(t)

		tests := []struct {
			input    float64
			kv       []string
			expected bool
		}{
			{43.345, []string{"max:43.346"}, true},
			{43.346, []string{"max:43.346"}, true},
			{43.347, []string{"max:43.346"}, false},
			{43.345, []string{"min:43.346"}, false},
			{43.346, []string{"min:43.346"}, true},
			{43.347, []string{"min:43.346"}, true},
			{832.428, []string{"max:1111.1111", "min:333.33"}, true},
			{333.323, []string{"max:1111.1111", "min:333.33"}, false},
			{1111.1112, []string{"max:1111.1111", "min:333.33"}, false},
		}

		for _, test := range tests {
			assert.Equal(test.expected, floatValidator(test.input, test.kv))
		}
	})

	t.Run("String Validator", func(t *testing.T) {
		assert := assert.New(t)

		tests := []struct {
			input    string
			kv       []string
			expected bool
		}{
			{"ABCDEFGHIJKLMNOPQ", []string{"max:18"}, true},
			{"ABCDEFGHIJKLMNOPQR", []string{"max:18"}, true},
			{"ABCDEFGHIJKLMNOPQRS", []string{"max:18"}, false},
			{"ABCDEFGHIJKLMNOPQ", []string{"min:18"}, false},
			{"ABCDEFGHIJKLMNOPQR", []string{"min:18"}, true},
			{"ABCDEFGHIJKLMNOPQRS", []string{"min:18"}, true},
			{"", []string{"max:12", "min:1"}, false},
			{"A", []string{"max:12", "min:1"}, true},
			{"ABCDEFGHIJKLMNOPQRS", []string{"max:12", "min:1"}, false},
		}

		for _, test := range tests {
			assert.Equal(test.expected, stringValidator(test.input, test.kv))
		}
	})
}

func TestValidate(t *testing.T) {
	type house struct {
		Cost    int    `validator:"number,max:10000000"`
		Address string `validator:"string,max:100"`
	}

	type asset struct {
		Cash   int `validator:"number,max:1000000"`
		Credit int `validator:"number,max:5000000"`
		Houses []house
	}

	type person struct {
		Name    string `validator:"string,min:6,max:20"`
		Phone   string `validator:"string,regex:^09[0-9]{2}-[0-9]{6}$"`
		Address string `validator:"string,max:70"`
		Age     int    `validator:"number,min:1,max:120"`
		Height  int    `validator:"number,min:20,max:250"`
		Asset   asset
		Email   string
	}

	t.Run("Person Layer", func(t *testing.T) {
		assert := assert.New(t)

		var tests = []struct {
			input    person
			expected bool
			errMsg   string
		}{
			{
				person{
					Name:    "Septem",
					Phone:   "0909-112331",
					Address: "1633  Hillview Street",
					Email:   "p1@gmail.com",
					Age:     199,
					Height:  180,
				}, false, `Age over 120, shouldn't be passed`,
			}, {
				person{
					Name:    "Nicole",
					Phone:   "0912-432122A",
					Address: "4051  Ethels Lane",
					Email:   "p2@gmail.com",
					Age:     18,
					Height:  162,
				}, false, `Phone format not correct, shouldn't be passed`,
			}, {
				person{
					Name:    "Pneumonoultramicroscopicsilicovolcanoconiosis",
					Phone:   "0912-445567",
					Address: "3603  Hood Avenue",
					Email:   "p3@yahoo.com.tw",
					Age:     43,
					Height:  200,
				}, false, `Name length over maximum value, shouldn't be passed`,
			}, {
				person{
					Name:    "Asolia",
					Phone:   "0912-345678",
					Address: "4018  Felosa Drive",
					Email:   "p4@gmail.com",
					Age:     50,
					Height:  180,
				}, true, `All informations are correct, should be passed`,
			},
		}

		for _, test := range tests {
			assert.Equal(Validate(test.input), test.expected, test.errMsg)
		}
	})

	t.Run("Asset Layer", func(t *testing.T) {
		assert := assert.New(t)

		var tests = []struct {
			input    person
			expected bool
			errMsg   string
		}{
			{
				person{
					Name:    "Septem",
					Phone:   "0909-090909",
					Address: "4872  Cameron Road",
					Email:   "p1@gmail.com",
					Age:     43,
					Height:  118,
					Asset: asset{
						Cash:   2000000,
						Credit: 200000,
					},
				}, false, `Cash of Asset over maximum value, shouldn't be passed`,
			}, {
				person{
					Name:    "Nicole",
					Phone:   "0912-321321",
					Address: "4798  Hickory Ridge Drive",
					Email:   "p2@gmail.com",
					Age:     99,
					Height:  179,
					Asset: asset{
						Cash:   123123,
						Credit: 10000000,
					},
				}, false, `Credit of Asset over maximum value, shouldn't be passed`,
			}, {
				person{
					Name:    "PassMan",
					Phone:   "0912-978978",
					Address: "2599  Sycamore Fork Road",
					Email:   "p3@gmail.com",
					Age:     25,
					Height:  190,
					Asset: asset{
						Cash:   100000,
						Credit: 500000,
					},
				}, true, `All of Asset informations are correct, should be passed`,
			},
		}

		for _, test := range tests {
			assert.Equal(Validate(test.input), test.expected, test.errMsg)
		}
	})

	t.Run("House Layer", func(t *testing.T) {
		assert := assert.New(t)

		var tests = []struct {
			input    person
			expected bool
			errMsg   string
		}{
			{
				person{
					Name:    "KerKer",
					Phone:   "0909-123123",
					Address: "3449  Jim Rosa Lane",
					Email:   "p1@gmail.com",
					Age:     38,
					Height:  168,
					Asset: asset{
						Cash:   100000,
						Credit: 300000,
						Houses: []house{
							house{
								Cost:    150000000,
								Address: "4382  Creekside Lane",
							},
							house{
								Cost:    100001,
								Address: "2587  Froe Street",
							},
						},
					},
				}, false, `House 1 of Asset cost over maximum value, shouldn't be passed`,
			}, {
				person{
					Name:    "Porter",
					Phone:   "0988-448792",
					Address: "3743  Clarence Court",
					Email:   "p2@gmail.com",
					Age:     65,
					Height:  200,
					Asset: asset{
						Cash:   200000,
						Credit: 500000,
						Houses: []house{
							house{
								Cost:    500000,
								Address: "3567  Valley Lane",
							},
							house{
								Cost:    600000,
								Address: "4891  Smith Road",
							},
						},
					},
				}, true, `All of Asset informations are correct, shouldn't be passed`,
			},
		}

		for _, test := range tests {
			assert.Equal(Validate(test.input), test.expected, test.errMsg)
		}
	})
}
