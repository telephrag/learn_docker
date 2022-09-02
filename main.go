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
		fmt.Fprintf(rw, "мон тынад секретъёсыд тодисько: %s\n", s)
	})

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, syscall.SIGTERM, syscall.SIGINT)
	<-interupt

	// Turns out that `docker stop %container%` sends SIGTERM.
	// If it doesn't work `docker kill --signal=SIGKILL` is executed implicitly.
	// In result line bellow will execute after container is stopped.
	// In fact killing with any signal specified inside `Notify()` will result in execution
	// of the line bellow.
	// NOTE: To see execute container via `docker run` cause idk how to execute restart
	//		 in non-headless mode.

	log.Println("gracefull shutdown after container stops")
}
