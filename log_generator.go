package main

import (
    "fmt"
    "net/http"
    // "html/template"
  	"io/ioutil"
    "log"
    "os"
    "time"
    "math/rand"
    "io"
)

func outputhandler(w http.ResponseWriter, r *http.Request) {
  file, err := ioutil.ReadFile("apache_log.txt") // For read access.
  if err != nil {
    http.Redirect(w, r, "/coffee", http.StatusFound)
    } else {
      log.Println(string(file[:]))
      fmt.Fprintf(w, string(file[:]))
    }
}

func random_with_ticker_handler(w http.ResponseWriter, r*http.Request) {
  log_files := [4]string{"apache_log.txt", "stack_trace.txt", "json_blob.txt" , "logseeds.txt"}
  filename := log_files[rand.Intn(len(log_files))]
  random_logfile, err := ioutil.ReadFile(filename)
  ticker := time.NewTicker(time.Millisecond * 10000)
  if err != nil {
      fmt.Fprintf(w, "Denied.")
  } else {
    random_tick_logfile, log_err := os.OpenFile("random_tick_logs.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if log_err != nil {
      log.Fatal("error opening log file: %v", err)
    } else {
      go func() {
          for t := range ticker.C {
              fmt.Println("Tick at", t)
              log.SetOutput(random_tick_logfile)
              log.Println(string(random_logfile[:]))
          }
      }()
      time.Sleep(time.Millisecond * 100000)
      ticker.Stop()
      fmt.Println("Ticker stopped")
      }
    }
}



func random_loghandler(w http.ResponseWriter, r*http.Request) {
  log_files := [4]string{"apache_log.txt", "stack_trace.txt", "json_blob.txt" , "logseeds.txt"}
  filename := log_files[rand.Intn(len(log_files))]
  random_logfile, err := ioutil.ReadFile(filename)
  if err != nil {
      fmt.Fprintf(w, "Denied.")
  } else {
      logfile, log_err := os.OpenFile("randomized_logs.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
      if log_err != nil {
        log.Fatal("error opening log file: %v", err)
      }
      defer logfile.Close()
      fmt.Fprintf(w, "Logs written to randomized_logs.txt from %s", filename)
      log.SetOutput(io.MultiWriter(logfile, os.Stdout, os.Stderr))
      log.Println(string(random_logfile[:]))
  }
}

func tickerloghandler(w http.ResponseWriter, r*http.Request) {
  ticker := time.NewTicker(time.Millisecond * 10000)
  f, _ := os.OpenFile("tickerlogfile", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    go func() {
        for t := range ticker.C {
            fmt.Println("Tick at", t)
            log.SetOutput(f)
            log.Println("logs logs better > bad == good")
        }
    }()
    time.Sleep(time.Millisecond * 100000)
    ticker.Stop()
    fmt.Println("Ticker stopped")
    log.SetOutput(f)
    log.Println("all done for now...")
}

func main() {
    http.HandleFunc("/tickfile", random_with_ticker_handler)
    http.HandleFunc("/output", outputhandler)
    http.HandleFunc("/tickerlog", tickerloghandler)
    http.HandleFunc("/random", random_loghandler)
    http.ListenAndServe(":8080", nil)
}
