# jsontree
Simple JSON Tree Key Traversing And Searching

If given a simple JSON tree containing nothing but ids as keys and children as arrays of keys like so:

```javascript
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
```

jsontree package gives you tools to Traverse and Search keys to return dot-notated paths fro use in other packages like gjson.
