package jsontree

import (
	"testing"

	"github.com/bmiles-development/gjson"
	"github.com/bmilesp/sjson"
	"github.com/stretchr/testify/assert"
)

func TestGetTopmostAncestorId(t *testing.T) {
	id, _ := GetTopmostAncestorId(testJsonTree)
	assert.Equal(t, `a`, id)

	value := gjson.Get(testJsonTree, "a.0.b.1.d.0")
	id, _ = GetTopmostAncestorId(value.String())
	assert.Equal(t, `e`, id)
}

func TestGetNextYoungerSiblingId(t *testing.T) {
	id, _ := GetNextYoungerSiblingId(testJsonTree, "g")
	assert.Equal(t, `h`, id)

	id, _ = GetNextYoungerSiblingId(testJsonTree, "i")
	assert.Equal(t, ``, id)

	id, _ = GetNextYoungerSiblingId(testJsonTree, "a")
	assert.Equal(t, ``, id)

	id, _ = GetNextYoungerSiblingId(testJsonTree, "b")
	assert.Equal(t, `m`, id)

	id, _ = GetNextYoungerSiblingId(testJsonTree, "m")
	assert.Equal(t, `n`, id)

	id, _ = GetNextYoungerSiblingId(testJsonTree, "k")
	assert.Equal(t, `l`, id)
}

func TestHasChildren(t *testing.T) {
	id, _ := HasChildren(testJsonTree, "b")
	assert.True(t, id)

	id, _ = HasChildren(testJsonTree, "c")
	assert.False(t, id)

	id, _ = HasChildren(testJsonTree, "d")
	assert.True(t, id)

	id, _ = HasChildren(testJsonTree, "e")
	assert.True(t, id)

	id, _ = HasChildren(testJsonTree, "n")
	assert.False(t, id)
}

func TestGetFirstChildId(t *testing.T) {
	id, _ := GetFirstChildId(testJsonTree, "b")
	assert.Equal(t, `c`, id, "they should be equal")

	id, _ = GetFirstChildId(testJsonTree, "c")
	assert.Equal(t, ``, id, "they should be equal")

	id, _ = GetFirstChildId(testJsonTree, "d")
	assert.Equal(t, `e`, id, "they should be equal")

	id, _ = GetFirstChildId(testJsonTree, "e")
	assert.Equal(t, `f`, id, "they should be equal")

	id, _ = GetFirstChildId(testJsonTree, "n")
	assert.Equal(t, ``, id, "they should be equal")
}

func TestGetElderSiblingId(t *testing.T) {
	id, _ := GetElderSiblingId(testJsonTree, "b")
	assert.Equal(t, ``, id, "they should be equal")

	id, _ = GetElderSiblingId(testJsonTree, "j")
	assert.Equal(t, ``, id, "they should be equal")

	id, _ = GetElderSiblingId(testJsonTree, "m")
	assert.Equal(t, `b`, id, "they should be equal")

	id, _ = GetElderSiblingId(testJsonTree, "n")
	assert.Equal(t, `m`, id, "they should be equal")

	id, _ = GetElderSiblingId(testJsonTree, "l")
	assert.Equal(t, `k`, id, "they should be equal")

	id, _ = GetElderSiblingId(testJsonTree, "a")
	assert.Equal(t, ``, id, "they should be equal")

	id, _ = GetElderSiblingId(testJsonTree, "g")
	assert.Equal(t, `f`, id, "they should be equal")
}

func TestGetIdFromPAth(t *testing.T) {
	id := getIdfromPath("a.0.b")
	assert.Equal(t, id, `b`, "they should be equal")

	id = getIdfromPath("a")
	assert.Equal(t, id, `a`, "they should be equal")

	id = getIdfromPath("a.0.b.1.d.0.e.3.i.0.j")
	assert.Equal(t, id, `j`, "they should be equal")
}

func TestGetNumericArrayKeyFromPath(t *testing.T) {
	id, _ := getNumericArrayKeyFromPath("a.0.b")
	assert.Equal(t, id, 0, "they should be equal")

	id, _ = getNumericArrayKeyFromPath("a.2.n")
	assert.Equal(t, id, 2, "they should be equal")

	id, _ = getNumericArrayKeyFromPath("a.0.b.1.d.0.e.3.i.0.j")
	assert.Equal(t, id, 0, "they should be equal")

	id, _ = getNumericArrayKeyFromPath("a.0.b.1.d.0.e.3.i")
	assert.Equal(t, id, 3, "they should be equal")
}

func TestGetYoungerSiblingsIds(t *testing.T) {
	res, _ := GetYoungerSiblingsIds(testJsonTree, "g")
	expected := []string{"h", "i"}
	assert.Equal(t, expected, res, "they should be equal")

	res, _ = GetYoungerSiblingsIds(testJsonTree, "f")
	expected = []string{"g", "h", "i"}
	assert.Equal(t, expected, res, "they should be equal")

	res, _ = GetYoungerSiblingsIds(testJsonTree, "b")
	expected = []string{"m", "n"}
	assert.Equal(t, expected, res, "they should be equal")

	res, _ = GetYoungerSiblingsIds(testJsonTree, "m")
	expected = []string{"n"}
	assert.Equal(t, expected, res, "they should be equal")

	res, _ = GetYoungerSiblingsIds(testJsonTree, "n")
	var expected2 []string
	assert.Equal(t, expected2, res, "they should be equal")

	res, _ = GetYoungerSiblingsIds(testJsonTree, "a")
	assert.Equal(t, expected2, res, "they should be equal")
}

func TestFlattenJson(t *testing.T) {

	res, _ := flattenJson(testJsonTreeSimple)
	assert.Equal(t, res, `{"a.0.b":null}`, "they should be equal")

	res, _ = flattenJson(testJsonTree)
	assert.Equal(t, res, bigResult1, "they should be equal")
}

func TestGetPathFromId(t *testing.T) {

	data, _ := flattenJson(testJsonTreeSimple)
	res, _ := getPathFromId(data, "b")
	assert.Equal(t, res, "a.0.b", "they should be equal")

	res, _ = getPathFromId(data, "a")
	assert.Equal(t, res, "a", "they should be equal")

	data, _ = flattenJson(testJsonTree)
	res, _ = getPathFromId(data, "a")
	assert.Equal(t, res, "a", "they should be equal")

	res, _ = getPathFromId(data, "j")
	assert.Equal(t, res, "a.0.b.1.d.0.e.3.i.0.j", "they should be equal")

	res, _ = getPathFromId(data, "d")
	assert.Equal(t, res, "a.0.b.1.d", "they should be equal")

	res, _ = getPathFromId(data, "m")
	assert.Equal(t, res, "a.1.m", "they should be equal")
}

func TestGetDescendants(t *testing.T) {
	res, _ := getDescendants(testJsonTreeSimple, "a")
	assert.Equal(t, res, `[{"b" : []}]`, "they should be equal")

	res, _ = getDescendants(testJsonTreeSimple, "b")
	assert.Equal(t, res, `[]`, "they should be equal")

	res, _ = getDescendants(testJsonTree, "e")
	assert.Equal(t, res, `[{"f":[]},{"g":[]},{"h":[]},{"i":[{"j":[]},{"k":[]},{"l":[]}]}]`, "they should be equal")

	res, _ = getDescendants(testJsonTree, "i")
	assert.Equal(t, res, `[{"j":[]},{"k":[]},{"l":[]}]`, "they should be equal")

	res, _ = getDescendants(testJsonTree, "a")
	assert.Equal(t, res, `[{"b":[{"c":[]},{"d":[{"e":[{"f":[]},{"g":[]},{"h":[]},{"i":[{"j":[]},{"k":[]},{"l":[]}]}]}]}]},{"m":[]},{"n":[]}]`, "they should be equal")

	res, _ = getDescendants(testJsonTree, "l")
	assert.Equal(t, res, `[]`, "they should be equal")
}

func TestGetParentPath(t *testing.T) {
	flatTree, _ := flattenJson(testJsonTreeSimple)
	res, _ := getParentPath(flatTree, "b")
	assert.Equal(t, res, "a", "they should be equal")

	res, _ = getParentPath(flatTree, "a")
	assert.Equal(t, res, "", "they should be equal")

	flatTree, _ = flattenJson(testJsonTree)
	res, _ = getParentPath(flatTree, "j")
	assert.Equal(t, res, `a.0.b.1.d.0.e.3.i`, "they should be equal")

	res, _ = getParentPath(flatTree, "a")
	assert.Equal(t, res, "", "they should be equal")

	res, _ = getParentPath(flatTree, "m")
	assert.Equal(t, res, "a", "they should be equal")
}

func TestGetElementNumberPath(t *testing.T) {

	res, _ := getElementNumberPath("a.0.b")
	assert.Equal(t, res, `a.0`, "they should be equal")

	res, _ = getElementNumberPath("a.0.b.1.d.0.e.3.i")
	assert.Equal(t, res, `a.0.b.1.d.0.e.3`, "they should be equal")

	res, _ = getElementNumberPath("a")
	assert.Equal(t, res, ``, "they should be equal")
}

func TestGetSiblingNumericPathsById(t *testing.T) {
	res, _ := getSiblingNumericPathsById(testJsonTree, "g")
	expected := []string{"a.0.b.1.d.0.e.0", "a.0.b.1.d.0.e.2", "a.0.b.1.d.0.e.3"}
	assert.Equal(t, res, expected, "they should be equal")

	res, _ = getSiblingNumericPathsById(testJsonTree, "h")
	expected = []string{"a.0.b.1.d.0.e.0", "a.0.b.1.d.0.e.1", "a.0.b.1.d.0.e.3"}
	assert.Equal(t, res, expected, "they should be equal")

	res, _ = getSiblingNumericPathsById(testJsonTree, "i")
	expected = []string{"a.0.b.1.d.0.e.0", "a.0.b.1.d.0.e.1", "a.0.b.1.d.0.e.2"}
	assert.Equal(t, res, expected, "they should be equal")

	res, _ = getSiblingNumericPathsById(testJsonTree, "l")
	expected = []string{"a.0.b.1.d.0.e.3.i.0", "a.0.b.1.d.0.e.3.i.1"}
	assert.Equal(t, res, expected, "they should be equal")

	res, _ = getSiblingNumericPathsById(testJsonTree, "a")
	expected = nil
	assert.Equal(t, res, expected, "they should be equal")

	res, _ = getSiblingNumericPathsById(testJsonTree, "n")
	expected = []string{"a.0", "a.1"}
	assert.Equal(t, res, expected, "they should be equal")
}

func TestGetDistinctFromNumericPath(t *testing.T) {
	res, _ := getDistinctFromNumericPath(testJsonTree, "a.0.b.1.d.0.e.0")
	assert.Equal(t, res, "a.0.b.1.d.0.e.0.f", "they should be equal")

	res, _ = getDistinctFromNumericPath(testJsonTree, "a.0.b.1")
	assert.Equal(t, res, "a.0.b.1.d", "they should be equal")

	res, _ = getDistinctFromNumericPath(testJsonTree, "a.0")
	assert.Equal(t, res, "a.0.b", "they should be equal")

	res, _ = getDistinctFromNumericPath(testJsonTree, "a.0.b.1.d.0.e.3.i.0")
	assert.Equal(t, res, "a.0.b.1.d.0.e.3.i.0.j", "they should be equal")
}

func TestGetAllSiblingsIds(t *testing.T) {
	res, _ := GetAllSiblingsIds(testJsonTree, "g")
	expected := []string{"f", "h", "i"}
	assert.Equal(t, res, expected, "they should be equal")

	res, _ = GetAllSiblingsIds(testJsonTree, "a")
	expected = []string(nil)
	assert.Equal(t, res, expected, "they should be equal")

	res, _ = GetAllSiblingsIds(testJsonTree, "b")
	expected = []string{"m", "n"}
	assert.Equal(t, res, expected, "they should be equal")

	res, _ = GetAllSiblingsIds(testJsonTree, "l")
	expected = []string{"j", "k"}
	assert.Equal(t, res, expected, "they should be equal")

	res, _ = GetAllSiblingsIds(testJsonTree, "c")
	expected = []string{"d"}
	assert.Equal(t, res, expected, "they should be equal")
}

func TestGetDescendantsIds(t *testing.T) {
	res, _ := GetDescendantsIds(testJsonTree, "g", false)
	expected := []string(nil)
	assert.Equal(t, expected, res, "they should be equal")

	res, _ = GetDescendantsIds(testJsonTree, "i", false)
	expected = []string{"j", "k", "l"}
	assert.Equal(t, expected, res, "they should be equal")

	res, _ = GetDescendantsIds(testJsonTree, "d", false)
	expected = []string{"e", "f", "g", "h", "i", "j", "k", "l"}
	assert.Equal(t, expected, res, "they should be equal")

	res, _ = GetDescendantsIds(testJsonTree, "n", false)
	expected = []string(nil)
	assert.Equal(t, expected, res, "they should be equal")

	res, _ = GetDescendantsIds(testJsonTree, "a", false)
	expected = []string{"b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n"}
	assert.Equal(t, expected, res, "they should be equal")

	res, _ = GetDescendantsIds(testJsonTree, "e", true)
	expected = []string{"f", "g", "h", "i"}
	assert.Equal(t, expected, res, "they should be equal")

}

func TestGetParentId(t *testing.T) {
	res, _ := GetParentId(testJsonTree, "i")
	expected := "e"
	assert.Equal(t, expected, res, "they should be equal")

	res, _ = GetParentId(testJsonTree, "h")
	expected = "e"
	assert.Equal(t, expected, res, "they should be equal")

	res, _ = GetParentId(testJsonTree, "f")
	expected = "e"
	assert.Equal(t, expected, res, "they should be equal")

	res, _ = GetParentId(testJsonTree, "b")
	expected = "a"
	assert.Equal(t, expected, res, "they should be equal")

	res, _ = GetParentId(testJsonTree, "m")
	expected = "a"
	assert.Equal(t, expected, res, "they should be equal")

	res, _ = GetParentId(testJsonTree, "a")
	expected = ""
	assert.Equal(t, expected, res, "they should be equal")

}

func TestIsFirstChild(t *testing.T) {
	res, _ := IsFirstChild(testJsonTree, "a")
	assert.True(t, res)

	res, _ = IsFirstChild(testJsonTree, "b")
	assert.True(t, res)

	res, _ = IsFirstChild(testJsonTree, "f")
	assert.True(t, res)

	res, _ = IsFirstChild(testJsonTree, "j")
	assert.True(t, res)

	res, _ = IsFirstChild(testJsonTree, "i")
	assert.False(t, res)

	res, _ = IsFirstChild(testJsonTree, "m")
	assert.False(t, res)

	res, _ = IsFirstChild(testJsonTree, "n")
	assert.False(t, res)

	res, _ = IsFirstChild(testJsonTree, "l")
	assert.False(t, res)
}

func TestIsLastChild(t *testing.T) {
	res, _ := IsLastChild(testJsonTree, "h")
	assert.False(t, res)

	res, _ = IsLastChild(testJsonTree, "a")
	assert.True(t, res)

	res, _ = IsLastChild(testJsonTree, "b")
	assert.False(t, res)

	res, _ = IsLastChild(testJsonTree, "j")
	assert.False(t, res)

	res, _ = IsLastChild(testJsonTree, "i")
	assert.True(t, res)

	res, _ = IsLastChild(testJsonTree, "m")
	assert.False(t, res)

	res, _ = IsLastChild(testJsonTree, "n")
	assert.True(t, res)

	res, _ = IsLastChild(testJsonTree, "l")
	assert.True(t, res)

}

func TestAddNextToLeafById(t *testing.T) {
	res, _ := AddNextToLeafById(testJsonTree, `h`, `{"w": [{"y":[]}]}`, "before")
	assert.Equal(t, `{"a":[{"b":[{"c":[]},{"d":[{"e":[{"f":[]},{"g":[]},{"w":[{"y":[]}]},{"h":[]},{"i":[{"j":[]},{"k":[]},{"l":[]}]}]}]}]},{"m":[]},{"n":[]}]}`, res, "they should be equal")

	res, _ = AddNextToLeafById(testJsonTree, `f`, `{"w": [{"y":[]}]}`, "before")
	assert.Equal(t, `{"a":[{"b":[{"c":[]},{"d":[{"e":[{"w":[{"y":[]}]},{"f":[]},{"g":[]},{"h":[]},{"i":[{"j":[]},{"k":[]},{"l":[]}]}]}]}]},{"m":[]},{"n":[]}]}`, res, "they should be equal")

	res, _ = AddNextToLeafById(testJsonTree, `l`, `{"w": [{"y":[]}]}`, "before")
	assert.Equal(t, `{"a":[{"b":[{"c":[]},{"d":[{"e":[{"f":[]},{"g":[]},{"h":[]},{"i":[{"j":[]},{"k":[]},{"w":[{"y":[]}]},{"l":[]}]}]}]}]},{"m":[]},{"n":[]}]}`, res, "they should be equal")

	res, _ = AddNextToLeafById(testJsonTree, `l`, `{"w": [{"y":[]}]}`, "after")
	assert.Equal(t, `{"a":[{"b":[{"c":[]},{"d":[{"e":[{"f":[]},{"g":[]},{"h":[]},{"i":[{"j":[]},{"k":[]},{"l":[]},{"w":[{"y":[]}]}]}]}]}]},{"m":[]},{"n":[]}]}`, res, "they should be equal")

	res, _ = AddNextToLeafById(testJsonTree, `b`, `{"xxx":[]}`, "after")
	assert.Equal(t, `{"a":[{"b":[{"c":[]},{"d":[{"e":[{"f":[]},{"g":[]},{"h":[]},{"i":[{"j":[]},{"k":[]},{"l":[]}]}]}]}]},{"xxx":[]},{"m":[]},{"n":[]}]}`, res, "they should be equal")

}

func TestAddIntoLeafById(t *testing.T) {
	res, _ := AddIntoLeafById(testJsonTree, `h`, `{"w": [{"y":[]}]}`, "insideBeginning")
	assert.Equal(t, `{"a":[{"b":[{"c":[]},{"d":[{"e":[{"f":[]},{"g":[]},{"h":[{"w":[{"y":[]}]}]},{"i":[{"j":[]},{"k":[]},{"l":[]}]}]}]}]},{"m":[]},{"n":[]}]}`, res, "they should be equal")

	res, _ = AddIntoLeafById(testJsonTree, `h`, `{"w": [{"y":[]}]}`, "insideEnd")
	assert.Equal(t, `{"a":[{"b":[{"c":[]},{"d":[{"e":[{"f":[]},{"g":[]},{"h":[{"w":[{"y":[]}]}]},{"i":[{"j":[]},{"k":[]},{"l":[]}]}]}]}]},{"m":[]},{"n":[]}]}`, res, "they should be equal")

	res, _ = AddIntoLeafById(testJsonTree, `f`, `{"w": [{"y":[]}]}`, "insideEnd")
	assert.Equal(t, `{"a":[{"b":[{"c":[]},{"d":[{"e":[{"f":[{"w":[{"y":[]}]}]},{"g":[]},{"h":[]},{"i":[{"j":[]},{"k":[]},{"l":[]}]}]}]}]},{"m":[]},{"n":[]}]}`, res, "they should be equal")

	res, _ = AddIntoLeafById(testJsonTree, `a`, `{"w": [{"y":[]}]}`, "insideBeginning")
	assert.Equal(t, `{"a":[{"w":[{"y":[]}]},{"b":[{"c":[]},{"d":[{"e":[{"f":[]},{"g":[]},{"h":[]},{"i":[{"j":[]},{"k":[]},{"l":[]}]}]}]}]},{"m":[]},{"n":[]}]}`, res, "they should be equal")

	res, _ = AddIntoLeafById(testJsonTree, `l`, `{"w": [{"y":[]}]}`, "insideEnd")
	assert.Equal(t, `{"a":[{"b":[{"c":[]},{"d":[{"e":[{"f":[]},{"g":[]},{"h":[]},{"i":[{"j":[]},{"k":[]},{"l":[{"w":[{"y":[]}]}]}]}]}]}]},{"m":[]},{"n":[]}]}`, res, "they should be equal")

	res, _ = AddIntoLeafById(testJsonTree, `b`, `{"xxx":[]}`, "insideEnd")
	assert.Equal(t, `{"a":[{"b":[{"c":[]},{"d":[{"e":[{"f":[]},{"g":[]},{"h":[]},{"i":[{"j":[]},{"k":[]},{"l":[]}]}]}]},{"xxx":[]}]},{"m":[]},{"n":[]}]}`, res, "they should be equal")

}

func TestAddJsonPieces(t *testing.T) {
	//get parent path from id:
	flatTree, _ := flattenJson(testJsonTree)

	parentPath, _ := getParentPath(flatTree, "i")
	assert.Equal(t, `a.0.b.1.d.0.e`, parentPath, "they should be equal")

	//find total siblings
	totalChildren := gjson.Get(testJsonTree, parentPath+".#")
	assert.Equal(t, float64(4), totalChildren.Num, "they should be equal")

	branch := gjson.Parse(testJsonTree).Get("a.0.b.1.d.0.e")
	assert.Equal(t, `[{"f":[]},{"g":[]},{"h":[]},{"i":[{"j":[]},{"k":[]},{"l":[]}]}]`, branch.Raw, "they should be equal")

	//loop through branch and append where needs be

	value, _ := sjson.Set(testJsonTree, "a.0.b.3", gjson.Parse(`{"w": [{"y":[]}]}`).Value().(map[string]interface{}))
	assert.Equal(t, value, `{"a":[{"b":[{"c":[]},{"d":[{"e":[{"f":[]},{"g":[]},{"h":[]},{"i":[{"j":[]},{"k":[]},{"l":[]}]}]}]},null,{"w":[{"y":[]}]}]},{"m":[]},{"n":[]}]}`, "they should be equal")

	/*
		value, _ = sjson.Set(`[]`, "0", gjson.Parse(`{"w": [{"y":[]}]}`).Value().(map[string]interface{}))
		value, _ = sjson.Set(value, "1", gjson.Parse(`{"w": [{"y":[]}]}`).Value().(map[string]interface{}))
	*/

	//test delete
	deleted, _ := sjson.Delete(testJsonTree, "a.0.b.1")
	assert.Equal(t, `{"a":[{"b":[{"c":[]}]},{"m":[]},{"n":[]}]}`, deleted, "they should be equal")

}

func TestRemoveById(t *testing.T) {

	res, _ := RemoveById(testJsonTreeSimple, "b")
	//log.Println(res)
	assert.Equal(t, `{"a":[]}`, res)

	res, _ = RemoveById(testJsonTree, "b")
	//log.Println(res)
	assert.Equal(t, `{"a":[{"m":[]},{"n":[]}]}`, res)

	res, _ = RemoveById(testJsonTree, `i`)
	//log.Println(res)
	assert.Equal(t, `{"a":[{"b":[{"c":[]},{"d":[{"e":[{"f":[]},{"g":[]},{"h":[]}]}]}]},{"m":[]},{"n":[]}]}`, res)

	_, err := RemoveById(testJsonTree, `a`)
	assert.Error(t, err)

	_, err = RemoveById(testJsonTree, `asrdgb35h54`)
	assert.Error(t, err)

}

var testJsonTreeSimple = `{"a":[{"b" : []}]}`
var testJsonTree = `{"a":[{"b":[{"c":[]},{"d":[{"e":[{"f":[]},{"g":[]},{"h":[]},{"i":[{"j":[]},{"k":[]},{"l":[]}]}]}]}]},{"m":[]},{"n":[]}]}`

/*
{expected []string
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
