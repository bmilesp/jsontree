package jsontree

import (
	"encoding/json"
	"strings"

	flatten "github.com/bmilesp/gojsonexplode"
	gjson "github.com/tidwall/gjson"
)

var Delimiter = "."

func GetDescendants(json string, key string) (string, error) {
	flatTree, err := flattenJson(json)
	if err != nil {
		return "", err
	}
	parentPath, err := getPathFromKey(flatTree, key)
	if err != nil {
		return "", err
	}
	value := gjson.Get(json, parentPath)
	return value.String(), nil
}

func flattenJson(tree string) (string, error) {
	flat, err := flatten.Explodejsonstr(tree, Delimiter)
	if err != nil {
		return "", err
	}
	return flat, nil
}

func getPathFromKey(flatTree string, key string) (string, error) {
	m := make(map[string]string)
	err := json.Unmarshal([]byte(flatTree), &m)
	if err != nil {
		return "", err
	}
	for k, _ := range m {
		if strings.Contains(k, key) {
			splitKeys := strings.Split(k, Delimiter)
			var returnPath strings.Builder
			for _, v := range splitKeys {
				if v == key {
					returnPath.WriteString(key)
					return returnPath.String(), nil
				}
				returnPath.WriteString(v)
				returnPath.WriteString(Delimiter)
			}
		}
	}
	return "", nil
}

func getParentPath(flatTree string, key string) (string, error) {
	path, err := getPathFromKey(flatTree, key)
	if err != nil {
		return "", err
	}
	splitKeys := strings.Split(path, Delimiter)
	if len(splitKeys) >= 3 {
		splitKeys = splitKeys[:len(splitKeys)-2]
		parentPath := strings.Join(splitKeys[:], Delimiter)
		return parentPath, nil
	}

	return "", nil
}
