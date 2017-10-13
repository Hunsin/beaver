# beaver [![build status](https://api.travis-ci.org/Hunsin/beaver.svg?branch=master)](https://travis-ci.org/Hunsin/beaver) [![GoDoc](https://godoc.org/github.com/Hunsin/beaver?status.svg)](https://godoc.org/github.com/Hunsin/beaver)

A set of libraries which may be useful in different projects.
[Godoc](https://godoc.org/github.com/Hunsin/beaver)

## Install
`go get github.com/Hunsin/beaver`

## JSON
Example of reading/writing JSON file and GET/POST JSON from http services
```go
package main

import (
  "net/http"
  "github.com/Hunsin/beaver"
)

type example struct {
  Hello string `json:"hello"`
  Year  int    `json:"year"`
  IP    string `json:"ip"`
}

func main() {
  e := example{
    Hello: "Hello World!",
    Year:  2017,
  }

  // writing JSON file; ignore error
  beaver.JSON(&e).WriteFile("example.json")

  // reading JSON file; ignore error
  out := example{}
  js  := beaver.JSON(&out)
  js.Open("example.json")

  // now you have the values
  println(out.Hello) // Hello World!
  println(out.Year)  // 2017

  // serving JSON data; ignore error
  http.HandleFunc("/path", func(w http.ResponseWriter, r *http.Request) {
    js.Serve(w, http.StatusOK)
  })

  // get values from other service; ignore error
  // http header is nil in this case
  js.Get("http://ip.jsontest.com", nil)
  println(js.IP)

  // send data to another service
  // the second parameter is your custom http.Header
  err := js.Post("https://url.of.service", nil)
  if err != nil {
    println(err.Error)
  }
}
```