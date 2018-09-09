package jsontree

import (
	"testing"
  "github.com/stretchr/testify/assert"
)

func TestFlattenJson(t *testing.T) {

  res, _ := FlattenJson(testJsonTreeSimple)
  assert.Equal(t, res, result1, "they should be equal")

  res, _ = FlattenJson(testJsonTree)
  assert.Equal(t, res, result2, "they should be equal")
}


func TestGetPathFromKey(t *testing.T) {

  data, _ := FlattenJson(testJsonTreeSimple)
  res, _ := GetPathFromKey(data,"b");

  assert.Equal(t, res, "a.0.b", "they should be equal")

  res, _ = GetPathFromKey(data,"a");
  assert.Equal(t, res, "a", "they should be equal")

  data, _ = FlattenJson(testJsonTree)
  res, _ = GetPathFromKey(data,"a");
  assert.Equal(t, res, "a", "they should be equal")

  res, _ = GetPathFromKey(data,"j");
  assert.Equal(t, res, "a.0.b.1.d.0.e.3.i.0.j", "they should be equal")

  res, _ = GetPathFromKey(data,"d");
  assert.Equal(t, res, "a.0.b.1.d", "they should be equal")

  res, _ = GetPathFromKey(data,"m");
  assert.Equal(t, res, "a.1.m", "they should be equal")
}



var testJsonTreeSimple = `{"a":[{"b" : []}]}`
var testJsonTree = `{"a":[{"b":[{"c":[]},{"d":[{"e":[{"f":[]},{"g":[]},{"h":[]},{"i":[{"j":[]},{"k":[]},{"l":[]}]}]}]}]},{"m":[]},{"n":[]}]}`

var result1 = `{"a.0.b":null}`
var result2 = `{"a.0.b.0.c":null,"a.0.b.1.d.0.e.0.f":null,"a.0.b.1.d.0.e.1.g":null,"a.0.b.1.d.0.e.2.h":null,"a.0.b.1.d.0.e.3.i.0.j":null,"a.0.b.1.d.0.e.3.i.1.k":null,"a.0.b.1.d.0.e.3.i.2.l":null,"a.1.m":null,"a.2.n":null}`
