package main

import "io/ioutil"
import "fmt"
import "path"
import "net/url"
import "encoding/json"

type ErrorReporter struct {
  errors map[string]*err
}

func (r * ErrorReporter) Start () {
  r.errors = make(map[string]*err)
}

func (r * ErrorReporter) Found (u * url.URL) {
}

func (r * ErrorReporter) Finish (report string) {
  list := make([]*err, 0, 10000)

  for _, err := range r.errors {
    list = append(list, err)
  }

  json, err := json.MarshalIndent(list, "", "  ")
  if err != nil {
    fmt.Println(err)
    return
  }

  err = ioutil.WriteFile(path.Join(report, "errors.json"), json, 0644)
  if err != nil {
    fmt.Println(err)
    return
  }
}

func (r * ErrorReporter) Success (u * url.URL, _ uint) {
}

func (r * ErrorReporter) Ignored (u * url.URL, _ uint, reason interface{}) {
}

func (r * ErrorReporter) Error (u * url.URL, status uint, reason interface{}) {
  if reason != nil {
    reason = fmt.Sprintf("%v", reason)
  } else {
    reason = ""
  }
  r.errors[u.String()] = &err{ URL: u.String(), Status: status, Reason: reason.(string) }
}

type err struct {
  URL    string `json:"url"`
  Status uint   `json:"status,omitempty"`
  Reason string `json:"reason,omitempty"`
}
