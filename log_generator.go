package main

import (
    "net/http"
  	"io/ioutil"
    "log"
    "os"
    "time"
    "math/rand"
    "io"
    "strconv"
    "fmt"
)

var counter = 1

// run through log_seeds folder and grab all file names, push into a slice
func create_filename_slice() []string {
  files, err := ioutil.ReadDir("/log_seeds")
  check(err)
  filename_slice := make([]string, 0)
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

// append logfile with random log file contents
func append_logfile() {
  filename := get_random_logfile()
  random_logfile, err := ioutil.ReadFile("/log_seeds/" + filename)
  check(err)
  logfile, log_err := os.OpenFile("/var/log/all_logs.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
  check(log_err)
  defer logfile.Close()
  // log.SetOutput(io.MultiWriter(logfile, os.Stdout, os.Stderr))
  log.SetOutput(io.MultiWriter(logfile))
  log.Println("Log Entry #", counter, "\n", string(random_logfile[:]))
  counter += 1
}

// start a function that will push one log each (10 seconds) for 100 seconds
func random_with_ticker_handler(w http.ResponseWriter, r*http.Request) {
  ticker := time.NewTicker(time.Millisecond * 10000)
  go func() {
      for range ticker.C {
        append_logfile()
      }
  }()
  // time.Sleep(time.Millisecond * 100000)
  // ticker.Stop()
}

// add contents of one random file to all_logs.txt
func random_loghandler(w http.ResponseWriter, r*http.Request) {
  append_logfile()
  fmt.Fprintf(w, "Logs written to randomized_logs.txt.  Search for super-fantastic-amazing!\n")
  fmt.Fprintf(w, "Logs total: %d", counter)
}


// user can request a certain number of logs using batch?n=integer, will default to 20 if number not given
func batchhandler(w http.ResponseWriter, r*http.Request){
  str_number := r.URL.Query().Get("n")
  num, err := strconv.Atoi(str_number)
  if err != nil {
    num = 20
  }
  for i := 0; i < num ; i++ {
    append_logfile()
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
