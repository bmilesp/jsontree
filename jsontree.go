package jsontree

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	flatten "github.com/bmilesp/gojsonexplode"
	gjson "github.com/tidwall/gjson"
)

var Delimiter = "."

func GetParentId(jsonTree string, key string) (string, error) {
	flatTree, err := flattenJson(jsonTree)
	if err != nil {
		return "", err
	}
	parentPath, err := getParentPath(flatTree, key)
	if err != nil {
		return "", err
	}
	parentId := getIdfromPath(parentPath)
	return parentId, err
}

func GetDescendantsIds(jsonTree string, key string) ([]string, error) {
	var ids []string
	var emptyIds []string
	descendantJsonTree, err := getDescendants(jsonTree, key)
	if err != nil {
		return ids, err
	}
	var m []interface{}
	err = json.Unmarshal([]byte(descendantJsonTree), &m)
	if err != nil {
		return ids, err
	}
	for _, v := range m {
		vals := v.(map[string]interface{})
		for k, _ := range vals {
			ids = append(ids, k)
			//check for more descendants
			descendantIds, err := GetDescendantsIds(jsonTree, k)
			if err != nil {
				return emptyIds, err
			}
			ids = append(ids, descendantIds...)
		}
	}
	return ids, err
}

func GetAllSiblingsIds(jsonTree string, id string) ([]string, error) {
	var ids []string
	paths, err := getSiblingNumericPathsById(jsonTree, id)
	if err != nil {
		return ids, err
	}
	//	var keys []string
	for _, v := range paths {
		value := gjson.Get(jsonTree, v)
		m := make(map[string]interface{})
		err := json.Unmarshal([]byte(value.String()), &m)
		if err != nil {
			return ids, err
		}
		for kk, _ := range m {
			ids = append(ids, kk)
		}
	}
	return ids, err
}

func GetFirstChildId(jsonTree string, key string) (string, error) {
	descendantIds, err := GetDescendantsIds(jsonTree, key)
	if err != nil {
		return "", err
	}
	for _, v := range descendantIds {
		return v, err
	}
	return "", errors.New("No children found.")
}

func HasChildren(jsonTree string, key string) (bool, error) {
	descendantIds, err := GetDescendantsIds(jsonTree, key)
	if err != nil {
		return false, err
	}
	if len(descendantIds) == 0 {
		return false, err
	}
	return true, err
}

func IsFirstChild(jsonTree string, key string) (bool, error) {
	flatTree, err := flattenJson(jsonTree)
	if err != nil {
		return false, err
	}
	path, err := getPathFromId(flatTree, key)
	if err != nil {
		return false, err
	}
	pathToElemNumber, err := getElementNumberPath(path)
	if err != nil {
		return false, err
	}
	//if "" then top element, so defaults to true
	if pathToElemNumber == "" {
		return true, err
	}
	arrayKey := pathToElemNumber[len(pathToElemNumber)-1:]
	if arrayKey == "0" {
		return true, err
	}
	return false, err
}

func GetYoungerSiblingsIds(jsonTree string, key string) ([]string, error) {
	var youngerSiblings []string
	flatTree, err := flattenJson(jsonTree)
	if err != nil {
		return youngerSiblings, err
	}
	path, err := getPathFromId(flatTree, key)
	if err != nil {
		return youngerSiblings, err
	}
	pathToElemNumber, err := getElementNumberPath(path)
	if err != nil {
		return youngerSiblings, err
	}
	//if "" then top element, so return empty slice
	if pathToElemNumber == "" {
		return youngerSiblings, err
	}
	childKeyString := pathToElemNumber[len(pathToElemNumber)-1:]
	if err != nil {
		return youngerSiblings, err
	}
	childKey, err := strconv.Atoi(childKeyString)
	if err != nil {
		return youngerSiblings, err
	}
	allSiblingNumericPaths, err := getSiblingNumericPathsById(jsonTree, key)
	if err != nil {
		return youngerSiblings, err
	}
	for _, currentNumericPath := range allSiblingNumericPaths {
		distinctPath, err := getDistinctFromNumericPath(jsonTree, currentNumericPath)
		if err != nil {
			return youngerSiblings, err
		}
		numericSiblingKey, err := getNumericArrayKeyFromPath(distinctPath)
		if err != nil {
			return youngerSiblings, err
		}
		if numericSiblingKey < childKey {
			continue
		}
		id := getIdfromPath(distinctPath)
		youngerSiblings = append(youngerSiblings, id)
	}
	return youngerSiblings, err
}

func GetElderSiblingId(jsonTree string, id string) (string, error) {
	elderSiblingId := ""
	parentId, err := GetParentId(jsonTree, id)
	if err != nil {
		return elderSiblingId, err
	}
	allChildren, err := GetAllSiblingsIds(jsonTree, id)
	if err != nil {
		return elderSiblingId, err
	}
	if allChildren == nil { // no elder siblings, first child
		return elderSiblingId, err
	}
	currentNumericKey, err := getNumericArrayKeyFromPath(id)
	for _, v := range allChildren {
		siblingNumericKey, err := getNumericArrayKeyFromPath(v)
		if err != nil {
			return "", err
		}
		if siblingNumericKey == currentNumericKey-1 {
			return getIdfromPath(v), err
		}
	}
	return parentId, err
}

func getIdfromPath(path string) string {
	popper := strings.Split(path, ".")
	if len(popper) >= 1 {
		return popper[len(popper)-1]
	}
	return ""
}

func getNumericArrayKeyFromPath(path string) (int, error) {
	var numericKey int
	pathToElemNumber, err := getElementNumberPath(path)
	if err != nil {
		return numericKey, err
	}
	//if "" then top element, so return empty slice
	if pathToElemNumber == "" {
		return numericKey, err
	}
	keyString := pathToElemNumber[len(pathToElemNumber)-1:]

	if err != nil {
		return numericKey, err
	}
	numericKey, err = strconv.Atoi(keyString)
	if err != nil {
		return numericKey, err
	}
	return numericKey, err
}

func getDistinctFromNumericPath(jsonTree string, path string) (string, error) {
	actualPath := ""
	value := gjson.Get(jsonTree, path)
	var m map[string]interface{}
	err := json.Unmarshal([]byte(value.String()), &m)
	if err != nil {
		return actualPath, err
	}
	for k, _ := range m {
		actualPath = path + Delimiter + k
		break
	}
	return actualPath, err
}

func getParentPath(flatTree string, id string) (string, error) {
	path, err := getPathFromId(flatTree, id)
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

func getSiblingNumericPathsById(json string, id string) ([]string, error) {
	var siblingPaths []string
	flatTree, err := flattenJson(json)
	if err != nil {
		return siblingPaths, err
	}
	path, err := getPathFromId(flatTree, id)
	if err != nil {
		return siblingPaths, err
	}
	pathToElemNumber, err := getElementNumberPath(path)
	if err != nil {
		return siblingPaths, err
	}
	parentPath, err := getParentPath(flatTree, id)
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

func getDescendants(jsonTree string, key string) (string, error) {
	flatTree, err := flattenJson(jsonTree)
	if err != nil {
		return "", err
	}
	parentPath, err := getPathFromId(flatTree, key)
	if err != nil {
		return "", err
	}
	value := gjson.Get(jsonTree, parentPath)
	return value.String(), nil
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

func getPathFromId(flatTree string, id string) (string, error) {
	m := make(map[string]string)
	err := json.Unmarshal([]byte(flatTree), &m)
	if err != nil {
		return "", err
	}
	for k, _ := range m {
		if strings.Contains(k, id) {
			splitKeys := strings.Split(k, Delimiter)
			var returnPath strings.Builder
			for _, v := range splitKeys {
				if v == id {
					returnPath.WriteString(id)
					return returnPath.String(), nil
				}
				returnPath.WriteString(v)
				returnPath.WriteString(Delimiter)
			}
		}
	}
	return "", nil
}
