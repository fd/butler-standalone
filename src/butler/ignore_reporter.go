package main

import "io/ioutil"
import "fmt"
import "path"
import "net/url"
import "encoding/json"

type IgnoreReporter struct {
  ignored map[string]*ignore
}

func (r * IgnoreReporter) Start () {
  r.ignored = make(map[string]*ignore)
}

func (r * IgnoreReporter) Finish (report string) {
  list := make([]*ignore, 0, 100)

  for _, ignore := range r.ignored {
    list = append(list, ignore)
  }

  json, err := json.MarshalIndent(list, "", "  ")
  if err != nil {
    fmt.Println(err)
    return
  }

  err = ioutil.WriteFile(path.Join(report, "ignored.json"), json, 0644)
  if err != nil {
    fmt.Println(err)
    return
  }
}

func (r * IgnoreReporter) Success (u * url.URL, _ uint) {
}

func (r * IgnoreReporter) Ignored (u * url.URL, status uint, reason interface{}) {
  if reason != nil {
    reason = fmt.Sprintf("%v", reason)
  } else {
    reason = ""
  }
  r.ignored[u.String()] = &ignore{ URL: u.String(), Status: status, Reason: reason.(string) }
}

func (r * IgnoreReporter) Error (u * url.URL, status uint, reason interface{}) {
}

type ignore struct {
  URL    string `json:"url"`
  Status uint   `json:"status,omitempty"`
  Reason string `json:"reason,omitempty"`
}
