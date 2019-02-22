package gospin

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestSpinner_Spin(t *testing.T) {
	rand.Seed(1)
	spinner := New(&Config{UseGlobalRand: true})

	simple := "The {slow|quick} {brown|blue and {red|yellow}} {fox|deer} {gracefully|} jumps over the {sleeping|lazy} dog"
	expected := "The slow blue and yellow deer jumps over the sleeping dog"
	got := spinner.Spin(simple)
	assert.Equal(t, expected, got, "should be equal")

	escaped := "The \\{escaped\\} {slow|quick} {brown|blue and {red|yellow}} {fox|deer} {gracefully|} jumps over the {sleeping|lazy} dog"
	expected = "The \\{escaped\\} slow brown deer gracefully jumps over the lazy dog"
	got = spinner.Spin(escaped)
	assert.Equal(t, expected, got, "should be equal")

	escapedDelimiters := "The \\{escaped\\} {slow|quick} {brown|blue and {red|yellow}} {fox|deer|escaped\\|pipe} {gracefully|} jumps over the {sleeping|lazy} dog"
	expected = "The \\{escaped\\} slow blue and red escaped\\pipe gracefully jumps over the lazy dog"
	got = spinner.Spin(escapedDelimiters)
	assert.Equal(t, expected, got, "should be equal")
}

func TestSpinner_Spin_WithCustomConfig(t *testing.T) {
	rand.Seed(1)
	spinner := New(&Config{
		StartChar:     "[",
		EndChar:       "]",
		EscapeChar:    "@",
		DelimiterChar: ";",
		UseGlobalRand: true,
	})

	simple := "The [slow;quick] [brown;blue and [red;yellow]] [fox;deer] [gracefully;] jumps over the [sleeping;lazy] dog"
	expected := "The slow blue and yellow deer jumps over the sleeping dog"
	got := spinner.Spin(simple)
	assert.Equal(t, expected, got, "should be equal")

	escaped := "The @[escaped@] [slow;quick] [brown;blue and [red;yellow]] [fox;deer] [gracefully;] jumps over the [sleeping;lazy] dog"
	expected = "The @[escaped@] slow brown deer gracefully jumps over the lazy dog"
	got = spinner.Spin(escaped)
	assert.Equal(t, expected, got, "should be equal")

	escapedDelimiters := "The @[escaped@] [slow;quick] [brown;blue and [red;yellow]] [fox;deer;escaped@;pipe] [gracefully;] jumps over the [sleeping;lazy] dog"
	expected = "The @[escaped@] slow blue and red escaped@pipe gracefully jumps over the lazy dog"
	got = spinner.Spin(escapedDelimiters)
	assert.Equal(t, expected, got, "should be equal")
}

func TestSpinner_Spin_WithNilConfig(t *testing.T) {
	spinner := New(nil)
	assert.Equal(t, defaultEscapeChar, spinner.EscapeChar)
	assert.Equal(t, defaultStartChar, spinner.StartChar)
	assert.Equal(t, defaultEndChar, spinner.EndChar)
	assert.Equal(t, defaultDelimiterChar, spinner.DelimiterChar)
}

func TestSpinner_SpinN(t *testing.T) {
	rand.Seed(1)
	spinner := New(&Config{UseGlobalRand: true})

	simple := "The {slow|quick} {brown|blue and {red|yellow}} {fox|deer} {gracefully|} jumps over the {sleeping|lazy} dog"
	expected := "The slow blue and yellow deer jumps over the sleeping dog"
	got := spinner.SpinN(simple, 100)
	assert.Equal(t, expected, got[0], "should be equal")
	assert.Len(t, got, 100, "should be equal")
}

func TestSpinner_Spin_NoDuplicates(t *testing.T) {
	spinner := New(nil)

	simple := "The {slow|quick} {brown|blue and {red|yellow}} {fox|deer} {gracefully|} jumps over the {sleeping|lazy} dog"
	got := spinner.SpinN(simple, 100)
	assert.Len(t, got, 100, "should be equal")

	ok := false

Exit:
	for _, str1 := range got {
		for _, str2 := range got {
			if str1 != str2 {
				ok = true
				continue Exit
			}
		}
	}

	if !ok {
		assert.Fail(t, "no unique values found")
	}
}
