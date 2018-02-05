# beaver [![build status](https://api.travis-ci.org/Hunsin/beaver.svg?branch=master)](https://travis-ci.org/Hunsin/beaver) [![GoDoc](https://godoc.org/github.com/Hunsin/beaver?status.svg)](https://godoc.org/github.com/Hunsin/beaver)

A set of functions which may be useful when dealing with JSON or logger.  
For more information please visit [Godoc](https://godoc.org/github.com/Hunsin/beaver).

## Install
`go get github.com/Hunsin/beaver`

## JSON
Example of reading/writing JSON file and GET/POST JSON from http services.
```go
package main

import (
  "net/http"
  bv "github.com/Hunsin/beaver"
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
  bv.JSON(&e).WriteFile("example.json")

  // reading JSON file; ignore error
  out := example{}
  js  := bv.JSON(&out)
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

## Logger
Logger wraps log.Logger with additional option to set log levels.
Let's see the example:
```go
package main

import (
  "log"
  "os"
  bv "github.com/Hunsin/beaver"
)

func main() {
  f, _ := os.Create("file_name.log")
  defer f.Close()

  // you can write usual and error log in different io.Writer
  bv.LogOutput(f, os.Stderr)

  // you can decide what level of logs should write, by default all
  // there are five levels available: Fatal, Error, Warn, Info & Debug
  bv.LogLevel(bv.Lerror | bv.Linfo)

  bv.Info("Hello World!") // 2018/02/06 00:31:28 INFO : Hello World!
  
  // you can define your log tag style
  t := bv.LTag{"| Fatal | ", "| Error | ", "| Warn | ", "| Info | ", "| Debug | "}

  // you can chain functions in configuration
  l := bv.NewLogger().Output(f, os.Stderr).Tags(t).Flags(log.Lshortfile)
  l.Error("Hello again!!") // main.go:27: | Error | Hello again!!
}
```