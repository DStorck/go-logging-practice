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
)

var counter = 1

// run through log_seeds folder and grab all file names, push into a slice
func create_filename_slice() []string {
  files, err := ioutil.ReadDir("log_seeds")
  filename_slice := make([]string, 0)
	check(err)
	for _, file := range files {
    filename_slice = append(filename_slice, file.Name())
	}
  return filename_slice
}

// pick a random log file name
func get_random_logfile() string {
  all_filenames := create_filename_slice()
  filename := all_filenames[rand.Intn(len(all_filenames))]
  return filename
}

// start a function that will push one log each (10 seconds) for 100 seconds
func random_with_ticker_handler(w http.ResponseWriter, r*http.Request) {
  ticker := time.NewTicker(time.Millisecond * 10000)
  random_tick_logfile, log_err := os.OpenFile("random_tick_logs.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
  check(log_err)
  go func() {
      for t := range ticker.C {
        filename := get_random_logfile()
        random_logfile, err := ioutil.ReadFile("log_seeds/" + filename)
        check(err)
        fmt.Println("Tick at", t)
        log.SetOutput(random_tick_logfile)
        log.Println("Log Entry #", counter, "\n", string(random_logfile[:]))
        fmt.Println("counter: ", counter)
        counter += 1
      }
  }()
  time.Sleep(time.Millisecond * 100000)
  ticker.Stop()
}

// function that pushes one random log to a file
func random_loghandler(w http.ResponseWriter, r*http.Request) {
  filename := get_random_logfile()
  random_logfile, err := ioutil.ReadFile("log_seeds/" + filename)
  check(err)
  logfile, log_err := os.OpenFile("randomized_logs.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
  check(log_err)
  fmt.Fprintf(w, "Logs written to randomized_logs.txt from %s", filename)
  log.SetOutput(io.MultiWriter(logfile, os.Stdout, os.Stderr))
  log.Println("Log Entry #", counter, "\n", string(random_logfile[:]))
  fmt.Println("counter: " , counter)
  counter += 1
}

// user can request a certain number of logs using batch?n=integer, will default to 20 if number not given
func batchhandler(w http.ResponseWriter, r*http.Request){
  logfile, log_err := os.OpenFile("batchlogs.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
  str_number := r.URL.Query().Get("n")
  num, err := strconv.Atoi(str_number)
  check(log_err)
  if err != nil {
    num = 20
  }
  for i := 0; i < num ; i++ {
    filename := get_random_logfile()
    random_logfile, err := ioutil.ReadFile("log_seeds/" + filename)
    check(err)
    log.SetOutput(io.MultiWriter(logfile, os.Stdout, os.Stderr))
    log.Println("Log Entry #", counter, "\n", string(random_logfile[:]))
    fmt.Println("counter: ", counter)
    counter += 1
  }
}

func check(err error) {
  if err != nil {
        panic(err)
  }
}

func main() {
    http.HandleFunc("/batch", batchhandler) // expecting ?n= number
    http.HandleFunc("/ticker", random_with_ticker_handler)
    http.HandleFunc("/random", random_loghandler)
    http.ListenAndServe(":8080", nil)
}
