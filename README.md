# beaver [![build status](https://api.travis-ci.org/Hunsin/beaver.svg?branch=master)](https://travis-ci.org/Hunsin/beaver) [![GoDoc](https://godoc.org/github.com/Hunsin/beaver?status.svg)](https://godoc.org/github.com/Hunsin/beaver)

A set of libraries which may be useful in different projects.
[Godoc](https://godoc.org/github.com/Hunsin/beaver)

## Install
`go get github.com/Hunsin/beaver`

## JSON
Example of reading/writing JSON file
```go
package main

import (
  "net/http"
  "github.com/Hunsin/beaver"
)

type example struct {
  Hello string `json:"hello"`
  Year  string `json:"year"`
}

type ip struct {
  IP string `json:ip`
}

func main() {
  e := example{
    Hello: "Hello World!",
    Year:  2017,
  }

  // example of writing JSON file; ignore error
  beaver.JSON(&e).WriteFile("path/to/file.json")

  // example of reading JSON file; ignore error
  out := example{}
  js  := beaver.JSON(&out)
  js.Open("path/to/file.json")

  // now you have the values
  println(out.Hello) // Hello World!
  println(out.Year)  // 2017

  // example of serving JSON data; ignore error
  http.HandleFunc("/path", func(w http.ResponseWriter, r *http.Request) {
    js.Serve(w, http.StatusOK)
  })

  // example of get values from other service; ignore error
  addr := ip{}
  cli  := beaver.JSON(&addr)
  cli.Get("http://ip.jsontest.com")
  println(addr.IP) 

  // example of send data to another service
  // the second parameter is your custom http.Header
  err := cli.Post("https://url.of.service", nil)
  if err != nil {
    println(err.Error)
  }
}
```