package main

import "fmt"
import "net/url"

type StdoutReporter struct {
  known     uint
  processed uint
  succeeded uint
  ignored   uint
  errored   uint
}

func (r * StdoutReporter) Start () {
}

func (r * StdoutReporter) Finish (report string) {
  fmt.Println("\033[1000D\033[KDone!")
}

func (r * StdoutReporter) Found (u * url.URL) {
  r.known += 1
}

func (r * StdoutReporter) Success (u * url.URL, status uint) {
  r.processed += 1
  r.succeeded += 1
  fmt.Printf("\033[1000D\033[K(%d/%d) \033[32m[%d]\033[0m: %.80s", r.processed, r.known, status, u)
}

func (r * StdoutReporter) Ignored (u * url.URL, status uint, reason interface{}) {
  r.processed += 1
  r.ignored   += 1
  if status > 0 {
    if reason != nil {
      fmt.Printf("\033[1000D\033[K(%d/%d) \033[33m[%d]\033[0m: %.80s (%v)", r.processed, r.known, status, u, reason)
    } else {
      fmt.Printf("\033[1000D\033[K(%d/%d) \033[33m[%d]\033[0m: %.80s", r.processed, r.known, status, u)
    }
  } else {
    if reason != nil {
      fmt.Printf("\033[1000D\033[K(%d/%d) \033[33m[IGN]\033[0m: %.80s (%v)", r.processed, r.known, u, reason)
    } else {
      fmt.Printf("\033[1000D\033[K(%d/%d) \033[33m[IGN]\033[0m: %.80s", r.processed, r.known, u)
    }
  }
}

func (r * StdoutReporter) Error (u * url.URL, status uint, reason interface{}) {
  r.processed += 1
  r.errored   += 1
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
