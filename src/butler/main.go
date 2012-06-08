package main

import "sync"
import "net/http"
import "net/url"
import "io/ioutil"
import "fmt"
import "html"
import "regexp"
import "strings"
import "os"
import "io"
import "bytes"
import "encoding/json"

import "math/rand"
import "hash/fnv"

import pqueue "github.com/nu7hatch/gopqueue"
import bloom  "github.com/dgryski/dgobloom"
import "haraway/common/options"

type Crawler struct {
  www         bool
  concurrency int
  domains     map[string]bool
  queue     * pqueue.Queue
  waiter      sync.WaitGroup
  known       bloom.BloomFilter

  report_dir  string
  reporters   []Reporter
}

type task struct {
  url * url.URL
}

func (t * task) Less (other interface{}) (bool) {
  return (len(t.url.String()) < len(other.(*task).url.String()))
}

func New(report_dir string)(c * Crawler, err error){

  salts_needed := bloom.SaltsRequired(100000, 0.001)
  salts := make([]uint32, salts_needed)
  for i := uint(0); i < salts_needed; i++ {
    salts[i] = rand.Uint32()
  }

  b := bloom.NewBloomFilter(100000, 0.001, fnv.New32(), salts)

  c = &Crawler{
    report_dir : report_dir,
    domains    : make(map[string]bool),
    queue      : pqueue.New(0),
    known      : b,
    reporters  : make([]Reporter, 0),
  }
  return
}

func (c * Crawler) RegisterReporter (reporter Reporter) {
  c.reporters = append(c.reporters, reporter)
}

func (c * Crawler) report_found (u * url.URL) {
  for _, reporter := range c.reporters {
    reporter.Found(u)
  }
}

func (c * Crawler) report_success (u * url.URL, status uint) {
  for _, reporter := range c.reporters {
    reporter.Success(u, status)
  }
}

func (c * Crawler) report_ignored (u * url.URL, status uint, reason interface{}) {
  for _, reporter := range c.reporters {
    reporter.Ignored(u, status, reason)
  }
}

func (c * Crawler) report_error (u * url.URL, status uint, reason interface{}) {
  for _, reporter := range c.reporters {
    reporter.Error(u, status, reason)
  }
}

func (c * Crawler) allow (domain string) {
  c.domains[domain] = true
}

func (c * Crawler) enqueue (link * url.URL, base * url.URL) {
  if base != nil {
    link = base.ResolveReference(link)
    link.Fragment = ""
  }

  if link.Path == "" {
    link.Path = "/"
  }

  link.Fragment = ""

  if link.Host != "" {
    link.Host = c.normalize_host(link.Host)
  }

  if c.known.Exists([]byte(link.String())) {
    return
  }

  c.known.Insert([]byte(link.String()))

  c.report_found(link)

  if link.Scheme == "http" {
    if c.domains[link.Host] {
      c.waiter.Add(1)
      c.queue.Enqueue(&task{url: link})
      return
    } else {
      c.report_ignored(link, 0, "external domain")
      return
    }
  } else {
    c.report_ignored(link, 0, "wrong scheme: "+link.Scheme)
    return
  }

  /*c.waiter.Add(1)*/
  /*c.queue <- u.String()*/
}

func (c * Crawler) Run () {
  os.MkdirAll(c.report_dir, 0755)

  for _, reporter := range c.reporters {
    reporter.Start()
  }

  for i := 0; i <= c.concurrency; i ++ {
    go func(){
      var buf bytes.Buffer
      for{
        t := c.queue.Dequeue()
        c.process_url(&buf, t.(*task).url)
        c.waiter.Done()
      }
    }()
  }

  c.waiter.Wait()

  for _, reporter := range c.reporters {
    reporter.Finish(c.report_dir)
  }
}

var pattern * regexp.Regexp

func (c * Crawler) process_url (buf * bytes.Buffer, page * url.URL) {
  defer buf.Reset()

  resp, err := http.Get(page.String())
  if err != nil {
    c.report_error(page, 0, err)
    return
  }
  defer resp.Body.Close()

  _, err = io.Copy(buf, resp.Body)
  if err != nil {
    c.report_error(page, 0, err)
    return
  }

  // check for redirects

  if resp.StatusCode != 200 {
    c.report_error(page, uint(resp.StatusCode), nil)
    return
  }

  if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
    c.report_ignored(page, 0, fmt.Sprintf("content-type: %v", resp.Header.Get("Content-Type")))
    return
  }

  links := pattern.FindAllSubmatch(buf.Bytes(), -1)
  for _, m := range links {
    link := string(m[1])

    link = html.UnescapeString(link)

    if strings.HasPrefix(link, "#") {
      continue
    }

    u, err := url.Parse(link)
    if err != nil { fmt.Printf("Invalid url: %s\n", link); continue }

    c.enqueue(u, page)
  }

  c.report_success(page, uint(resp.StatusCode))
}

func (c * Crawler) normalize_host (host string) (string) {
  if strings.HasPrefix(host, "www.") {
    if c.www {
      return host
    } else {
      return host[4:]
    }
  } else {
    if c.www {
      return "www." + host
    } else {
      return host
    }
  }
  return ""
}

func (c * Crawler) Load (path string) (err error) {
  var config   Config
  var u      * url.URL

  jsonBlob, err := ioutil.ReadFile(path)
  if err != nil { return }

  err = json.Unmarshal(jsonBlob, &config)
  if err != nil { return }

  c.www = config.Www

  if config.Concurrency <= 0 {
    config.Concurrency = 2
  }
  c.concurrency = config.Concurrency

  for _, domain := range config.Domains {
    domain = c.normalize_host(domain)
    c.allow(domain)

    u, err = url.Parse("http://"+domain+"/")
    if err != nil { return }

    c.enqueue(u, nil)
  }

  return
}

type Config struct {
  Concurrency int      `json:"concurrency"`
  Www         bool     `json:"www"`
  Domains     []string `json:"domains"`
}

func init() {
  var err error
  pattern, err = regexp.Compile("[<]a[^>]+href[=][\"']([^\"']+)[\"']")
  if err != nil { panic(err) }
}


const desc = `
butler - A simple sitemap generator and error reporter.
Usage: butler [--report=] [--config=]
--
report=  -r,--report=   Path to the report directory.
config=  -c,--config=   The path to the config file.
--
--
--

`

func main() {
  spec, err := options.New(desc)
  if err != nil { panic(err) }

  opts, err := spec.Parse(os.Args, os.Environ())

  if err != nil {
    spec.PrintUsageWithError(err)
  }

  if len(opts.Args) != 0 {
    spec.PrintUsageAndExit()
  }

  config := opts.Get("config")
  if config == "" { config = ".butler.json" }

  report := opts.Get("report")
  if report == "" { report = "report" }

  c, err := New(report)
  if err != nil { panic(err) }

  c.RegisterReporter(new(SitemapReporter))
  c.RegisterReporter(new(StdoutReporter))
  c.RegisterReporter(new(ErrorReporter))
  c.RegisterReporter(new(IgnoreReporter))

  err = c.Load(config)
  if err != nil { panic(err) }

  c.Run()
}
