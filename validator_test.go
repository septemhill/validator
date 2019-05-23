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
			{"090A-123432", []string{"regex:^09[0-9]{2}-[0-9]{6}$"}, false},
			{"0907123432", []string{"regex:^09[0-9]{2}-[0-9]{6}$"}, false},
			{"0907-123432", []string{"regex:^09[0-9]{2}-[0-9]{6}$"}, true},
		}

		for _, test := range tests {
			assert.Equal(test.expected, stringValidator(test.input, test.kv))
		}
	})
}

func TestNestedValidate(t *testing.T) {
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

func TestCollectionValidate(t *testing.T) {
	newFloat := func(v float64) *float64 {
		s := new(float64)
		*s = v
		return s
	}

	t.Run("Primitive Collection Validate", func(t *testing.T) {
		assert := assert.New(t)

		type collection struct {
			Array    [8]int     `validator:"min:10,max:20"`
			Slice    []string   `validator:"min:10,max:20"`
			SlicePtr []*float64 `validator:"min:10,max:20"`
		}

		tests := []struct {
			input    collection
			expected bool
		}{
			{
				collection{
					Array:    [8]int{1, 2, 3, 4, 5, 6, 7, 8},
					Slice:    []string{"Septem", "Nicole", "Asolia"},
					SlicePtr: []*float64{newFloat(12.343), newFloat(34.233), newFloat(97.8554)},
				}, false,
			}, {
				collection{
					Array:    [8]int{11, 12, 13, 14, 15, 16, 17, 18},
					Slice:    []string{"Septem1117", "Nicole0721", "Asolia0524"},
					SlicePtr: []*float64{newFloat(12.343), newFloat(11.233), newFloat(18.8554)},
				}, true,
			}, {
				collection{
					Array:    [8]int{19, 18, 17, 16, 15, 14, 13, 8},
					Slice:    []string{"AAAAAAAAAA", "BBBBBBBBBB", "CCCCCCCCCC", "DDDDDDDDDD", "EEEEEEEEEE"},
					SlicePtr: []*float64{newFloat(11.93), newFloat(13.3293), newFloat(13.5994)},
				}, false,
			},
		}

		for _, test := range tests {
			assert.Equal(test.expected, Validate(test.input))
		}
	})

	t.Run("Structure Collection Validate", func(t *testing.T) {
		assert := assert.New(t)

		type inner struct {
			Age    int     `validator:"min:18,max:120"`
			Height float64 `validator:"min:40,max:250"`
			Name   string  `validator:"min:6,max:20"`
		}

		type collection struct {
			Array    [2]inner
			Slice    []inner
			SlicePtr []*inner
		}

		tests := []struct {
			input    collection
			expected bool
		}{
			{
				collection{
					Array: [2]inner{
						inner{Age: 132, Height: 222, Name: "Nicker"},
						inner{Age: 33, Height: 123, Name: "Boasher"},
					},
					Slice: []inner{
						inner{Age: 100, Height: 188, Name: "Nicole"},
						inner{Age: 111, Height: 169, Name: "Buster"},
						inner{Age: 112, Height: 129, Name: "Nord"},
					},
					SlicePtr: []*inner{
						&inner{Age: 20, Height: 444, Name: "Dicker"},
						&inner{Age: 56, Height: 217, Name: "Zash"},
						&inner{Age: 34, Height: 211, Name: "Nasher"},
						&inner{Age: 67, Height: 142, Name: "Queen"},
					},
				}, false,
			}, {
				collection{
					Array: [2]inner{
						inner{Age: 56, Height: 123, Name: "Jackie"},
						inner{Age: 82, Height: 149, Name: "Steven"},
					},
					Slice: []inner{
						inner{Age: 19, Height: 139, Name: "Neo"},
					},
					SlicePtr: []*inner{
						&inner{Age: 28, Height: 187, Name: "Docker"},
						&inner{Age: 66, Height: 169, Name: "Chambers"},
					},
				}, false,
			}, {
				collection{
					Array: [2]inner{
						inner{Age: 20, Height: 120, Name: "Septem"},
						inner{Age: 40, Height: 121, Name: "Asolia"},
					},
					Slice: []inner{
						inner{Age: 60, Height: 122, Name: "Michael"},
						inner{Age: 80, Height: 123, Name: "Joseph"},
						inner{Age: 100, Height: 124, Name: "Johnny"},
					},
					SlicePtr: []*inner{
						&inner{Age: 111, Height: 200, Name: "Austin"},
						&inner{Age: 112, Height: 126, Name: "Martin"},
						&inner{Age: 113, Height: 127, Name: "Wilson"},
					},
				}, true,
			},
		}

		for _, test := range tests {
			assert.Equal(test.expected, Validate(test.input))
		}
	})
}
