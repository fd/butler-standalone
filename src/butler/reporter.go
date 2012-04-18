package main

import "net/url"

type Reporter interface {
  Start()

  // Error(url *url.URL, status uint)
  Success(*url.URL, uint)

  // Error(u *url.URL, status uint, reason interface{})
  Error(*url.URL, uint, interface{})

  // Ignored(url *url.URL, status uint, reason interface{})
  Ignored(*url.URL, uint, interface{})

  // Finish(report string)
  Finish(string)
}
