package main

import (
    "fmt"
    "net/http"
    // "html/template"
  	"io/ioutil"
    "log"
    "os"
    "time"
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
    fmt.Fprintf(w, "Don't you just love %s?!", r.URL.Path[1:])
}

func outputhandler(w http.ResponseWriter, r *http.Request) {
  file, err := ioutil.ReadFile("file.go") // For read access.
     if err != nil {
	   log.Println("ha i made something write to the console!")
     http.Redirect(w, r, "/coffee", http.StatusFound)
     } else {
       fmt.Fprintf(w, string(file[:]))
     }
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

func logwriter(w http.ResponseWriter, r*http.Request) {
  f, err := os.OpenFile("testlogfile2", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
  if err != nil {
    log.Fatal("error opening file: %v", err)
  }
  defer f.Close()
  fmt.Fprintf(w, "Logs written to testlogfile.")
  log.SetOutput(f)
  log.Println("This is a test log entry")
}

func tickerhandler(w http.ResponseWriter, r*http.Request) {
  ticker := time.NewTicker(time.Millisecond * 1000)
  go func() {
      for t := range ticker.C {
          fmt.Println("Tick at", t)
      }
  }()
  time.Sleep(time.Millisecond * 10000)
  ticker.Stop()
  fmt.Println("Ticker stopped")
}

func tickerloghandler(w http.ResponseWriter, r*http.Request) {
  ticker := time.NewTicker(time.Millisecond * 10000)
  f, _ := os.OpenFile("tickerlogfile", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    go func() {
        for t := range ticker.C {
            fmt.Println("Tick at", t)
            log.SetOutput(f)
            log.Println("logggggggssssss")
        }
    }()
    time.Sleep(time.Millisecond * 100000)
    ticker.Stop()
    fmt.Println("Ticker stopped")
    log.SetOutput(f)
    log.Println("all done for now...")

}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func helloworld(t time.Time) {
	fmt.Printf("%v: Hello, World!\n", t)
}

func main() {
    doEvery(200*time.Millisecond, helloworld)
    http.HandleFunc("/", handler)
    http.HandleFunc("/logs", loghandler)
    http.HandleFunc("/output", outputhandler)
    http.HandleFunc("/writelogs", logwriter)
    http.HandleFunc("/ticker", tickerhandler)
    http.HandleFunc("/tickerlog", tickerloghandler)
    http.ListenAndServe(":8080", nil)
}
