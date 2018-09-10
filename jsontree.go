package jsontree

import (
	"encoding/json"
	"strconv"
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

func GetSiblingKeys(json string, key string) ([]string, error) {
	var siblingPaths []string
	flatTree, err := flattenJson(json)
	if err != nil {
		return siblingPaths, err
	}
	path, err := getPathFromKey(flatTree, key)
	if err != nil {
		return siblingPaths, err
	}
	pathToElemNumber, err := getElementNumberPath(path)
	if err != nil {
		return siblingPaths, err
	}
	parentPath, err := getParentPath(flatTree, key)
	if err != nil {
		return siblingPaths, err
	}
	n := 0
	currentPath := ""
	for {
		currentPath = parentPath + Delimiter + strconv.Itoa(n)
		value := gjson.Get(json, currentPath)
		if value.String() != "" {
			if pathToElemNumber != currentPath {
				siblingPaths = append(siblingPaths, currentPath)
			}
		} else {
			break
		}
		n++
	}
	return siblingPaths, nil
}

func getElementNumberPath(path string) (string, error) {
	splitKeys := strings.Split(path, Delimiter)
	if len(splitKeys) >= 3 {
		splitKeys = splitKeys[:len(splitKeys)-1]
		elemPath := strings.Join(splitKeys[:], Delimiter)
		return elemPath, nil
	}
	return "", nil
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
