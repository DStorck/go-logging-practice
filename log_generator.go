package main

import (
    "fmt"
    "net/http"
  	"io/ioutil"
    "log"
    "os"
    "time"
    "math/rand"
    "io"
    "strconv"
    // "net/url"
    // "net"
)

func create_filename_slice() []string {
  files, err := ioutil.ReadDir("log_seeds")
  filename_slice := make([]string, 0)
	check(err)

	for _, file := range files {
		fmt.Println(file.Name())
    filename_slice = append(filename_slice, file.Name())
	}

  fmt.Println(filename_slice)
  return filename_slice
}

func get_random_logfile() string {
  all_filenames := create_filename_slice()
  filename := all_filenames[rand.Intn(len(all_filenames))]
  return filename
}

func random_with_ticker_handler(w http.ResponseWriter, r*http.Request) {
  filename := get_random_logfile()
  random_logfile, err := ioutil.ReadFile("log_seeds/" + filename)
  ticker := time.NewTicker(time.Millisecond * 10000)
  check(err)
  random_tick_logfile, log_err := os.OpenFile("random_tick_logs.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
  check(log_err)
  go func() {
      for t := range ticker.C {
          fmt.Println("Tick at", t)
          log.SetOutput(random_tick_logfile)
          log.Println(string(random_logfile[:]))
      }
  }()
  time.Sleep(time.Millisecond * 100000)
  ticker.Stop()
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
    log.SetOutput(f)
    log.Println("all done for now...")
}

func random_loghandler(w http.ResponseWriter, r*http.Request) {
  filename := get_random_logfile()
  random_logfile, err := ioutil.ReadFile("log_seeds/" + filename)
  check(err)
  logfile, log_err := os.OpenFile("randomized_logs.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
  check(log_err)
  // defer logfile.Close()
  fmt.Fprintf(w, "Logs written to randomized_logs.txt from %s", filename)
  log.SetOutput(io.MultiWriter(logfile, os.Stdout, os.Stderr))
  log.Println(string(random_logfile[:]))
}

func batchhandler(w http.ResponseWriter, r*http.Request){
  logfile, log_err := os.OpenFile("batchlogs.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
  str_number := r.URL.Query().Get("n")
  fmt.Println(str_number)
  num, err := strconv.Atoi(str_number)
  check(log_err)
  if err != nil {
    num = 20
    fmt.Println(num)
  }
  for i := 0; i < num ; i++ {
    filename := get_random_logfile()
    random_logfile, err := ioutil.ReadFile("log_seeds/" + filename)
    check(err)
    log.SetOutput(io.MultiWriter(logfile, os.Stdout, os.Stderr))
    log.Println(string(random_logfile[:]))
  }
}

func check(err error) {
  if err != nil {
        panic(err)
  }
}

func main() {
    http.HandleFunc("/batch", batchhandler) // expecting ?n= number
    http.HandleFunc("/tickfile", random_with_ticker_handler)
    http.HandleFunc("/tickerlog", tickerloghandler)
    http.HandleFunc("/random", random_loghandler)
    http.ListenAndServe(":8080", nil)
}
