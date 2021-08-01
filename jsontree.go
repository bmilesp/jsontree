package jsontree

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"

	flatten "github.com/bmilesp/gojsonexplode"
	"github.com/bmilesp/sjson"
	gjson "github.com/tidwall/gjson"
)

var Delimiter = "."

func GetParentId(jsonTree string, key string) (string, error) {
	flatTree, err := flattenJson(jsonTree)
	if err != nil {
		return `{"error": "jsontree.GetParentId flattenJson Error"}`, err
	}
	parentPath, err := getParentPath(flatTree, key)
	if err != nil {
		return `{"error": "jsontree.GetParentId parentPath Error"}`, err
	}
	parentId := getIdfromPath(parentPath)
	return parentId, err
}

func GetDescendantsIds(jsonTree string, key string, childrenOnly bool) ([]string, error) {
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
			if !childrenOnly {
				descendantIds, err := GetDescendantsIds(jsonTree, k, false)
				if err != nil {
					return emptyIds, err
				}
				ids = append(ids, descendantIds...)
			}
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
	descendantIds, err := GetDescendantsIds(jsonTree, key, false)
	if err != nil {
		return "", err
	}
	for _, v := range descendantIds {
		return v, err
	}
	return "", errors.New("No children found.")
}

func HasChildren(jsonTree string, key string) (bool, error) {
	descendantIds, err := GetDescendantsIds(jsonTree, key, false)
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

func IsLastChild(jsonTree string, key string) (bool, error) {
	youngerSiblings, err := GetYoungerSiblingsIds(jsonTree, key)
	if err != nil {
		return false, err
	}
	if youngerSiblings == nil {
		return true, err
	} else {
		return false, err
	}
}

func GetNextYoungerSiblingId(jsonTree string, id string) (string, error) {
	youngerSiblingsIds, err := GetYoungerSiblingsIds(jsonTree, id)
	if err != nil {
		return "error", err
	}
	if youngerSiblingsIds == nil {
		return "", err
	}
	return youngerSiblingsIds[0], err
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
	allSiblings, err := GetAllSiblingsIds(jsonTree, id)
	if err != nil {
		return elderSiblingId, err
	}
	if allSiblings == nil { // no elder siblings, first child
		return elderSiblingId, err
	}
	flatTree, err := flattenJson(jsonTree)
	if err != nil {
		return elderSiblingId, err
	}
	currentPath, err := getPathFromId(flatTree, id)
	if err != nil {
		return elderSiblingId, err
	}
	currentNumericKey, err := getNumericArrayKeyFromPath(currentPath)
	for _, v := range allSiblings {
		siblingPath, err := getPathFromId(flatTree, v)
		if err != nil {
			return elderSiblingId, err
		}
		siblingNumericKey, err := getNumericArrayKeyFromPath(siblingPath)
		if err != nil {
			return "", err
		}
		if siblingNumericKey == currentNumericKey-1 {
			return v, err
		}
	}
	return "", err
}

func GetTopmostAncestorId(jsonTree string) (string, error) {

	firstKey := ""
	var m map[string]interface{}
	err := json.Unmarshal([]byte(jsonTree), &m)

	if err != nil {
		return firstKey, err
	}

	for key, _ := range m {
		firstKey = key
		break
	}

	return firstKey, err
}

/*
func Add(jsonTree string, directive string, id string, jsonObject ){

	switch directive {
	case "before":
		//test can add a jsonTree within a tree?
		siblingTree := gjson.Get();

		//count all siblings
		siblingIds := GetAllSiblingsIds(jsonTree, id)
		var newSiblings []string
		for k, v := range siblingIds {
				if(v == id){
					newSiblings = append(newSiblingIds, id)
				}
				newSiblingIds = append(newSiblingIds, id)
		}


		newPath, err := sjson.Set(jsonTree,  )
	case "after":

	case "insertFirst":

	case "insertLast":

}
*/
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

// will return empty string if top ancestor path is passed in
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
	return "", errors.New("no id/path found")
}

func AddNextToLeafById(jsonTree string, id string, insertBranch string, beforeAfter string) (string, error) {
	flatTree, err := flattenJson(jsonTree)
	if err != nil {
		return `{"error": "jsontree.addNextToLeafById - flattening tree"}`, err
	}
	idPath, err := getPathFromId(flatTree, id)
	if err != nil {
		return `{"error": "jsontree.addNextToLeafById - getPathFromId"}`, err
	}
	insertKey, err := getNumericArrayKeyFromPath(idPath)
	if err != nil {
		return `{"error": "jsontree.addNextToLeafById - getNumericArrayKeyFromPath"}`, err
	}
	parentPath, err := getParentPath(flatTree, id)
	if err != nil {
		return `{"error": "jsontree.addNextToLeafById - getParentPath"}`, err
	}
	totalChildren := gjson.Get(jsonTree, parentPath+".#")

	newBranchStr := `[]`
	j := 0
	for i := 0; i < int(totalChildren.Num); i++ {
		insertKeyFound := i == insertKey

		if insertKeyFound && beforeAfter == "before" {
			val := gjson.Parse(insertBranch).Value().(map[string]interface{})
			newBranchStr, err = sjson.Set(newBranchStr, strconv.Itoa(j), val)
			j++
			if err != nil {
				return `{"error": "sjson.Set 1"}`, err
			}
		}

		siblingVal := gjson.Get(jsonTree, parentPath+"."+strconv.Itoa(i))
		newBranchStr, err = sjson.Set(newBranchStr, strconv.Itoa(j), gjson.Parse(siblingVal.Raw).Value().(map[string]interface{}))
		j++
		if err != nil {
			log.Println("sjson.Set 2")
			return "", err
		}

		if insertKeyFound && beforeAfter == "after" {
			val := gjson.Parse(insertBranch).Value().(map[string]interface{})
			newBranchStr, err = sjson.Set(newBranchStr, strconv.Itoa(j), val)
			j++
			if err != nil {

				return `{"error": "sjson.Set 3"}`, err
			}
		}
	}

	newJsonTree, err := sjson.Set(jsonTree, parentPath, gjson.Parse(newBranchStr).Value())
	if err != nil {

		return `{"error": "sjson.Set 4"}`, err
	}

	//siblingBranch := gjson.Parse(siblingBranchStr)
	return newJsonTree, err
}

func AddIntoLeafById(jsonTree string, id string, insertBranch string, topBottom string) (string, error) {
	flatTree, err := flattenJson(jsonTree)
	if err != nil {
		return `{"error": "jsontree.addIntoLeafById - flattening tree"}`, err
	}
	idPath, err := getPathFromId(flatTree, id)
	if err != nil {
		return `{"error": "jsontree.addIntoLeafById - getPathFromId"}`, err
	}

	totalChildren := gjson.Get(jsonTree, idPath+".#")
	newBranchStr := ""

	if int(totalChildren.Num) == 0 {
		val := gjson.Parse(insertBranch).Value().(map[string]interface{})
		newBranchStr, err = sjson.Set(jsonTree, idPath+".0", val)
		if err != nil {
			return `{"error": "jsontree.addIntoLeafById - failed to insert into empty Leaf"}`, err
		}
	} else {
		if topBottom == "insideBeginning" {
			result := gjson.Get(jsonTree, idPath+".0")

			firstKey := ""
			result.ForEach(func(key, value gjson.Result) bool {
				firstKey = key.String()
				return true
			})

			newBranchStr, err = AddNextToLeafById(jsonTree, firstKey, insertBranch, "before")
			if err != nil {
				return `{"error": "jsontree.addIntoLeafById - failed to insert before Leaf"}`, err
			}
		} else if topBottom == "insideEnd" {
			lastChildKey := strconv.Itoa(int(totalChildren.Num) - 1)
			result := gjson.Get(jsonTree, idPath+"."+lastChildKey)
			lastKey := ""
			result.ForEach(func(key, value gjson.Result) bool {
				lastKey = key.String()
				return true
			})
			newBranchStr, err = AddNextToLeafById(jsonTree, lastKey, insertBranch, "after")

			if err != nil {
				return `{"error": "jsontree.addIntoLeafById - failed to insert after Leaf"}`, err
			}
		} else {
			return `{"error": "jsontree.addIntoLeafById - directive must be either insideTop or insideBottom"}`, err
		}
	}

	return newBranchStr, err

}

func RemoveById(jsonTree string, id string) (string, error) {
	var newBranchStr string
	var err error
	flatTree, err := flattenJson(jsonTree)
	if err != nil {
		return `{"error": "jsontree.RemoveById - flattening tree"}`, err
	}
	idPath, err := getPathFromId(flatTree, id)
	if err != nil {
		return `{"error": "jsontree.RemoveById - getPathFromId"}`, err
	}
	pathToElemNumber, err := getElementNumberPath(idPath)
	if err != nil {
		return `{"error": "jsontree.RemoveById - getElementNumberPath"}`, err
	}
	newBranchStr, err = sjson.Delete(jsonTree, pathToElemNumber)
	if err != nil {
		return `{"error": "jsontree.RemoveById - failed to delete leaf/branch - Id cannot be top-most ancestor or key is invalid"}`, err
	}
	return newBranchStr, err
}
