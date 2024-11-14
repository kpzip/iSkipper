package main

import (
  "net/http"
  "io"
  "log"
)


func main() {
  // get config
  // TODO
  // sign into iCringe
  resp, err := http.Get("https://student.iclicker.com/#/login")
  if err != nil {
    log.Fatalln(err)
  }

  body, err := io.ReadAll(resp.Body)
  if err != nil {
    log.Fatalln(err)
  }

  stringbody := string(body)
  log.Printf(stringbody)
}
