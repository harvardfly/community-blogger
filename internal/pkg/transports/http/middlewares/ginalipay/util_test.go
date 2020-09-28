package ginalipay

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMustLength(t *testing.T) {
	cases := []struct {
		Str      string
		Length   int
		Expected string
		Err      bool
	}{
		{
			Str:      "iOS商品",
			Length:   0,
			Expected: "",
			Err:      true,
		},
		{
			Str:      "iOS商品",
			Length:   4,
			Expected: "",
			Err:      true,
		},
		{
			Str:      "iOS商品",
			Length:   5,
			Expected: "iOS商品",
			Err:      false,
		},
		{
			Str:      "iOS商品",
			Length:   6,
			Expected: "iOS商品",
			Err:      false,
		},
	}

	for _, v := range cases {
		exp, err := MustLength(v.Str, v.Length)
		if v.Err {
			assert.NotEqual(t, nil, err)
		} else {
			assert.Equal(t, nil, err)
		}
		assert.Equal(t, v.Expected, exp)
	}
}

func TestTruncateString(t *testing.T) {
	cases := []struct {
		Str      string
		Length   int
		Expected string
	}{
		{
			Str:      "iOS商品",
			Length:   0,
			Expected: "",
		},
		{
			Str:      "iOS商品",
			Length:   1,
			Expected: "i",
		},
		{
			Str:      "iOS商品",
			Length:   2,
			Expected: "iO",
		},
		{
			Str:      "iOS商品",
			Length:   3,
			Expected: "iOS",
		},
		{
			Str:      "iOS商品",
			Length:   4,
			Expected: "iOS商",
		},
		{
			Str:      "iOS商品",
			Length:   5,
			Expected: "iOS商品",
		},
		{
			Str:      "iOS商品",
			Length:   6,
			Expected: "iOS商品",
		},
	}

	for _, v := range cases {
		res := TruncateString(v.Str, v.Length)
		assert.Equal(t, v.Expected, res)
	}
}
