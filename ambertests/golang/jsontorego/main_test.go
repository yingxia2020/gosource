package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertStringArray(t *testing.T) {
	t.Run("Array size is one", func(t *testing.T) {
		input := []string{"aa"}
		expectedResult := fmt.Sprintf("[%q]\n", "aa")
		result := convertStringArray(input)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("Array size greater than one", func(t *testing.T) {
		input := []string{"aa", "bb", "cc"}
		expectedResult := fmt.Sprintf("[%q, %q, %q]\n", "aa", "bb", "cc")
		result := convertStringArray(input)
		assert.Equal(t, expectedResult, result)
	})
}

func TestConvertArray(t *testing.T) {
	t.Run("Int type", func(t *testing.T) {
		input := []int{0, 2, 4}
		expectedResult := fmt.Sprintf("[%v, %v, %v]\n", 0, 2, 4)
		result := convertArray(intToInterfaceSlice(input))
		assert.Equal(t, expectedResult, result)
	})

	t.Run("String type", func(t *testing.T) {
		input := []string{"aa", "bb", "cc"}
		expectedResult := fmt.Sprintf("[%q, %q, %q]\n", "aa", "bb", "cc")
		result := convertArray(stringToInterfaceSlice(input))
		assert.Equal(t, expectedResult, result)
	})
}

func TestConvertIncludeValue(t *testing.T) {
	t.Run("Int type value", func(t *testing.T) {
		value := 100
		title := "foo"
		expectedResult := "fooValue := [100]\nincludes_value(fooValue, quote.foo)\n\n"
		result := convertIncludeValue(value, title)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("Bool type value", func(t *testing.T) {
		value := false
		title := "foo"
		expectedResult := "fooValue := [false]\nincludes_value(fooValue, quote.foo)\n\n"
		result := convertIncludeValue(value, title)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("String type value", func(t *testing.T) {
		value := "abc"
		title := "foo"
		expectedResult := fmt.Sprintf("fooValue := [%q]\nincludes_value(fooValue, quote.foo)\n\n", "abc")
		result := convertIncludeValue(value, title)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("Float type value", func(t *testing.T) {
		value := 0.35
		title := "foo"
		expectedResult := "fooValue := [0.35]\nincludes_value(fooValue, quote.foo)\n\n"
		result := convertIncludeValue(value, title)
		assert.Equal(t, expectedResult, result)
	})
}

func TestConvertIncludeArray(t *testing.T) {
	t.Run("string type", func(t *testing.T) {
		input := []string{"aa", "bb", "cc"}
		title := "foo"
		expectedResult := "fooValues := " + convertStringArray(input) + "includes_value(fooValues, quote.foo)\n\n"
		result := convertIncludeArray(stringToInterfaceSlice(input), title)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("int type", func(t *testing.T) {
		input := []int{0, 2, 4}
		title := "foo"
		expectedResult := "fooValues := " + convertArray(intToInterfaceSlice(input)) +
			"includes_value(fooValues, quote.foo)\n\n"
		result := convertIncludeArray(intToInterfaceSlice(input), title)
		assert.Equal(t, expectedResult, result)
	})
}

func TestConvertLimitString(t *testing.T) {
	t.Run("Smaller equal than", func(t *testing.T) {
		limitField := "<=5"
		title := "foo"
		expectedResult := "foo_limit := " + fmt.Sprintf("%q", "5") + "\nsmaller_equal_than(quote.foo, foo_limit)\n\n"
		result := convertLimitValue(limitField, title, true)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("Smaller than", func(t *testing.T) {
		limitField := "<5"
		title := "foo"
		expectedResult := "foo_limit := " + fmt.Sprintf("%q", "5") + "\nsmaller_than(quote.foo, foo_limit)\n\n"
		result := convertLimitValue(limitField, title, true)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("Bigger equal than", func(t *testing.T) {
		limitField := ">=5"
		title := "foo"
		expectedResult := "foo_limit := " + fmt.Sprintf("%q", "5") + "\nbigger_equal_than(quote.foo, foo_limit)\n\n"
		result := convertLimitValue(limitField, title, true)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("Bigger than", func(t *testing.T) {
		limitField := ">5"
		title := "foo"
		expectedResult := "foo_limit := " + fmt.Sprintf("%q", "5") + "\nbigger_than(quote.foo, foo_limit)\n\n"
		result := convertLimitValue(limitField, title, true)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("Equal to", func(t *testing.T) {
		limitField := "5"
		title := "foo"
		expectedResult := "foo_value := " + fmt.Sprintf("[%q]", "5") + "\nincludes_value(foo_value, quote.foo)\n\n"
		result := convertLimitValue(limitField, title, true)
		assert.Equal(t, expectedResult, result)
	})
}

func TestConvertLimitInt(t *testing.T) {
	t.Run("Smaller equal than", func(t *testing.T) {
		limitField := "<=5"
		title := "foo"
		expectedResult := "foo_limit := " + fmt.Sprintf("%v", "5") + "\nsmaller_equal_than(quote.foo, foo_limit)\n\n"
		result := convertLimitValue(limitField, title, false)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("Smaller than", func(t *testing.T) {
		limitField := "<5"
		title := "foo"
		expectedResult := "foo_limit := " + fmt.Sprintf("%v", "5") + "\nsmaller_than(quote.foo, foo_limit)\n\n"
		result := convertLimitValue(limitField, title, false)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("Bigger equal than", func(t *testing.T) {
		limitField := ">=5"
		title := "foo"
		expectedResult := "foo_limit := " + fmt.Sprintf("%v", "5") + "\nbigger_equal_than(quote.foo, foo_limit)\n\n"
		result := convertLimitValue(limitField, title, false)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("Bigger than", func(t *testing.T) {
		limitField := ">5"
		title := "foo"
		expectedResult := "foo_limit := " + fmt.Sprintf("%v", "5") + "\nbigger_than(quote.foo, foo_limit)\n\n"
		result := convertLimitValue(limitField, title, false)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("Equal to", func(t *testing.T) {
		limitField := "5"
		title := "foo"
		expectedResult := "foo_value := " + fmt.Sprintf("[%v]", "5") + "\nincludes_value(foo_value, quote.foo)\n\n"
		result := convertLimitValue(limitField, title, false)
		assert.Equal(t, expectedResult, result)
	})
}
