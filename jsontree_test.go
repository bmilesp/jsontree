package jsontree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlattenJson(t *testing.T) {

	res, _ := flattenJson(testJsonTreeSimple)
	assert.Equal(t, res, `{"a.0.b":null}`, "they should be equal")

	res, _ = flattenJson(testJsonTree)
	assert.Equal(t, res, bigResult1, "they should be equal")
}

func TestGetPathFromKey(t *testing.T) {

	data, _ := flattenJson(testJsonTreeSimple)
	res, _ := getPathFromKey(data, "b")
	assert.Equal(t, res, "a.0.b", "they should be equal")

	res, _ = getPathFromKey(data, "a")
	assert.Equal(t, res, "a", "they should be equal")

	data, _ = flattenJson(testJsonTree)
	res, _ = getPathFromKey(data, "a")
	assert.Equal(t, res, "a", "they should be equal")

	res, _ = getPathFromKey(data, "j")
	assert.Equal(t, res, "a.0.b.1.d.0.e.3.i.0.j", "they should be equal")

	res, _ = getPathFromKey(data, "d")
	assert.Equal(t, res, "a.0.b.1.d", "they should be equal")

	res, _ = getPathFromKey(data, "m")
	assert.Equal(t, res, "a.1.m", "they should be equal")
}

func TestGetDescendants(t *testing.T) {
	res, _ := GetDescendants(testJsonTreeSimple, "a")
	assert.Equal(t, res, `[{"b" : []}]`, "they should be equal")

	res, _ = GetDescendants(testJsonTreeSimple, "b")
	assert.Equal(t, res, `[]`, "they should be equal")

	res, _ = GetDescendants(testJsonTree, "e")
	assert.Equal(t, res, `[{"f":[]},{"g":[]},{"h":[]},{"i":[{"j":[]},{"k":[]},{"l":[]}]}]`, "they should be equal")

	res, _ = GetDescendants(testJsonTree, "i")
	assert.Equal(t, res, `[{"j":[]},{"k":[]},{"l":[]}]`, "they should be equal")

	res, _ = GetDescendants(testJsonTree, "a")
	assert.Equal(t, res, `[{"b":[{"c":[]},{"d":[{"e":[{"f":[]},{"g":[]},{"h":[]},{"i":[{"j":[]},{"k":[]},{"l":[]}]}]}]}]},{"m":[]},{"n":[]}]`, "they should be equal")

	res, _ = GetDescendants(testJsonTree, "l")
	assert.Equal(t, res, `[]`, "they should be equal")
}

var testJsonTreeSimple = `{"a":[{"b" : []}]}`
var testJsonTree = `{"a":[{"b":[{"c":[]},{"d":[{"e":[{"f":[]},{"g":[]},{"h":[]},{"i":[{"j":[]},{"k":[]},{"l":[]}]}]}]}]},{"m":[]},{"n":[]}]}`

/*
{
	"a": [{
		"b": [{
			"c": []
		}, {
			"d": [{
				"e": [{
					"f": []
				}, {
					"g": []
				}, {
					"h": []
				}, {
					"i": [{
						"j": []
					}, {
						"k": []
					}, {
						"l": []
					}]
				}]
			}]
		}]
	}, {
		"m": []
	}, {
		"n": []
	}]
}
*/

var bigResult1 = `{"a.0.b.0.c":null,"a.0.b.1.d.0.e.0.f":null,"a.0.b.1.d.0.e.1.g":null,"a.0.b.1.d.0.e.2.h":null,"a.0.b.1.d.0.e.3.i.0.j":null,"a.0.b.1.d.0.e.3.i.1.k":null,"a.0.b.1.d.0.e.3.i.2.l":null,"a.1.m":null,"a.2.n":null}`
