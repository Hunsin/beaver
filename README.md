# beaver

A set of libraries which may be useful in different projects.

`go get github.com/Hunsin/beaver`

## JSON
Example of reading/writing JSON file
```go
package main

import "github.com/Hunsin/beaver"

type example struct {
  Hello string `json:"hello"`
  Year  string `json:"year"`
}

func main() {
  e := example{
    Hello: "Hello World!",
    Year:  2017,
  }

  // example of writing JSON file
  err := beaver.WriteJSON("path/to/file.json", &e)
  if err != nil {
    println(err.Error())
  }

  // example of reading JSON file
  out := example{}
  err = beaver.OpenJSON("path/to/file.json", &out)
  if err != nil {
    println(err.Error())
  }

  // now you have the values
  println(out.Hello)
  println(out.Year)
}
```