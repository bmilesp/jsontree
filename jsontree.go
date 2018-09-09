package jsontree

import(
  "strings"
  flatten "github.com/bmilesp/gojsonexplode"
  "encoding/json"
)

var Delimiter = "."


func FlattenJson(tree string)(string, error){
  flat, err := flatten.Explodejsonstr(tree, Delimiter)
  if err != nil {
		return "", err
	}
  return flat, err
}


func GetPathFromKey(flatTree string,key string)(string, error){
  m := make(map[string]string)
  err := json.Unmarshal([]byte(flatTree), &m)

  if err != nil {
    return "", err
  }
  for k,_ := range m{
      if strings.Contains(k, key) {
  	    splitKeys := strings.Split(k, Delimiter)
        var returnPath strings.Builder
        for _,v := range splitKeys {
          if(v == key){
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
