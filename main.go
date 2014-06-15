package main

import (
	"net/http"
	"net/mail"
	"log"
	"fmt"
	"flag"
	"bufio"
	"os"
	"time"
)

const maxEmailSize = 100

var (
	port = 8080
	host = ""
	emailName = "email"
	filename = "emails.txt"
)


func init() {
	flag.IntVar(&port, "port", port, "web server port address")
	flag.StringVar(&host, "host", host, "web server host address")
	flag.StringVar(&emailName, "email", emailName, "email argument name")
	flag.StringVar(&filename, "filename", filename,  "file to save emails")
}

// writeLines writes the lines to the given file.
func writeLines(lines []string, path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(path)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func main() {
	http.HandleFunc("/save/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		email := r.FormValue(emailName)
		if email == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "%s param is requried", emailName)
			return
		}
		if len(email) > maxEmailSize {
			email = email[:maxEmailSize]
		}
		addr, err := mail.ParseAddress(email)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "%s param is malformed", emailName)
			return
		}
		writeLines([]string{fmt.Sprintf("%s|%s", addr.Address, time.Now())}, filename)

		fmt.Fprintf(w, "ok")
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil))

}
