package main

import (
    "fmt"
    "net/http"
    // "html/template"
  	"io/ioutil"
)

type Page struct {
    Title string
    Body  []byte
}

func loadPage(title string) (*Page, error) {
    filename := "logseeds.txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func loghandler(w http.ResponseWriter, r*http.Request) {
  if r.Method == "GET" {
       log_data, err := ioutil.ReadFile("logseeds.txt")
       if err != nil {
           fmt.Fprintf(w, "Denied.")
       } else {
           fmt.Fprintf(w, string(log_data[:]) + "\n")
       }
   }
}

func main() {
    http.HandleFunc("/index", handler)
    http.HandleFunc("/logs", loghandler)
    http.ListenAndServe(":8080", nil)
}
