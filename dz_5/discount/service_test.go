package discount

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatePercent(t *testing.T) {
	cases := []struct {
		name    string
		percent int
		want    error
	}{
		{"Zero", 0, nil},
		{"Five", 5, nil},
		{"Hundred", 100, nil},
		{"Negative", -1, ErrInvalidPercent},
		{"Overflow", 101, ErrInvalidPercent},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			err := ValidatePercent(c.percent)
			if c.want != nil {
				assert.ErrorIs(t, err, c.want)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestApply(t *testing.T) {
	cases := []struct {
		name       string
		price      int
		percent    int
		priceAfter int
		want       error
	}{
		{"30% discount", 1000, 30, 700, nil},
		{"0'%' discount", 1000, 0, 1000, nil},
		{"100% discount", 1000, 100, 0, nil},
		{"Negative price", -1, 30, 0, ErrNegativePrice},
		{"Bad discount", 1000, 101, 0, ErrInvalidPercent},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			got, err := Apply(c.price, c.percent)
			if c.want != nil {
				assert.ErrorIs(t, err, c.want)
			} else {
				require.NoError(t, err)
				assert.Equal(t, c.priceAfter, got)
			}
		})
	}
}

func TestFinalPrice(t *testing.T) {
	cases := []struct {
		name       string
		items      []Item
		totalPrice int
		want       error
	}{
		{"NIL", nil, 0, nil},
		{"Empty", []Item{}, 0, nil},
		{"Casual items", []Item{{100, 15}, {1000, 70}}, 385, nil},
		{"Full discount", []Item{{100, 100}}, 0, nil},
		{"Negative price", []Item{{-1, 10}}, 0, ErrNegativePrice},
		{"Negative percent", []Item{{100, -10}}, 0, ErrInvalidPercent},
		{"Overflow percent", []Item{{1000, 101}}, 0, ErrInvalidPercent},
		{"Invalid percent after valid", []Item{{1000, 10}, {150, -20}}, 0, ErrInvalidPercent},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			got, err := FinalPrice(c.items)
			if c.want != nil {
				assert.ErrorIs(t, err, c.want)
			} else {
				require.NoError(t, err)
				assert.Equal(t, c.totalPrice, got)
			}
		})
	}
}
