package main

import "os"
import "fmt"
import "path"
import "net/url"
import "container/list"

type SitemapReporter struct {
  urls map[string]*list.List
}

func (r * SitemapReporter) Start () {
  r.urls = make(map[string]*list.List)
}

func (r * SitemapReporter) Found (u * url.URL) {
}

func (r * SitemapReporter) Finish (report string) {
  os.Mkdir(path.Join(report, "sitemaps"), 0755)

  for host, urls := range r.urls {
    f, err := os.Create(path.Join(report, "sitemaps", host + ".xml"))
    defer f.Close()
    if err != nil {
      fmt.Println(err)
      return
    }

    f.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
    f.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")
    for e := urls.Front(); e != nil; e = e.Next() {
      f.WriteString("  <url>\n")
      f.WriteString("    <loc>"+e.Value.(*url.URL).String()+"</loc>\n")
      f.WriteString("  </url>\n")
    }
    f.WriteString("</urlset>\n")
  }
}

func (r * SitemapReporter) Success (u * url.URL, _ uint) {
  urls, present := r.urls[u.Host]
  if !present {
    urls = list.New()
    r.urls[u.Host] = urls
  }
  urls.PushBack(u)
}

func (r * SitemapReporter) Ignored (u * url.URL, _ uint, reason interface{}) {
}

func (r * SitemapReporter) Error (u * url.URL, _ uint, reason interface{}) {
}
