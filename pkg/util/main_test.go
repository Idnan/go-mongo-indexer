package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringInSlice(t *testing.T) {
	var actual bool

	actual = StringInSlice("plane", []string{"bus", "car", "bike"})
	assert.Equal(t, false, actual)

	actual = StringInSlice("bus", []string{"bus", "car", "bike"})
	assert.Equal(t, true, actual)
}

func TestJsonEncode(t *testing.T) {
	var actual string

	actual = JsonEncode([]string{"bus", "car", "bike"})
	assert.Equal(t, `["bus","car","bike"]`, actual)

	actual = JsonEncode(map[string]int{"bus": 200, "car": 300, "bike": 400})
	assert.Equal(t, `{"bike":400,"bus":200,"car":300}`, actual)
}


