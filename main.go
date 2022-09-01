package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/telephrag/errlist"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "chirtkem mudila")
		s, ok := syscall.Getenv("SECRET")
		if !ok {
			fmt.Fprintln(rw, errlist.New(fmt.Errorf("%d", http.StatusInternalServerError)))
		}
		fmt.Fprintf(rw, "мон тынад секретёсыд тодисько: %s\n", s)
	})

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, syscall.SIGTERM, syscall.SIGINT)
	<-interupt
}
