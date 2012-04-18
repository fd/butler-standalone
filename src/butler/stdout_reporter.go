package main

import "fmt"
import "net/url"

type StdoutReporter struct {
}

func (r * StdoutReporter) Start () {
}

func (r * StdoutReporter) Finish (report string) {
  fmt.Println("\033[1000D\033[KDone!")
}

func (r * StdoutReporter) Success (u * url.URL, status uint) {
  fmt.Printf("\033[1000D\033[K\033[32m[%d]\033[0m: %.80s", status, u)
}

func (r * StdoutReporter) Ignored (u * url.URL, status uint, reason interface{}) {
  if status > 0 {
    if reason != nil {
      fmt.Printf("\033[1000D\033[K\033[33m[%d]\033[0m: %.80s (%v)", status, u, reason)
    } else {
      fmt.Printf("\033[1000D\033[K\033[33m[%d]\033[0m: %.80s", status, u)
    }
  } else {
    if reason != nil {
      fmt.Printf("\033[1000D\033[K\033[33m[IGN]\033[0m: %.80s (%v)", u, reason)
    } else {
      fmt.Printf("\033[1000D\033[K\033[33m[IGN]\033[0m: %.80s", u)
    }
  }
}

func (r * StdoutReporter) Error (u * url.URL, status uint, reason interface{}) {
  if status > 0 {
    if reason != nil {
      fmt.Printf("\033[1000D\033[K\033[31m[%d]\033[0m: %.80s (%v)\n", status, u, reason)
    } else {
      fmt.Printf("\033[1000D\033[K\033[31m[%d]\033[0m: %.80s\n", status, u)
    }
  } else {
    if reason != nil {
      fmt.Printf("\033[1000D\033[K\033[31m[ERR]\033[0m: %.80s (%v)\n", u, reason)
    } else {
      fmt.Printf("\033[1000D\033[K\033[31m[ERR]\033[0m: %.80s\n", u)
    }
  }
}
