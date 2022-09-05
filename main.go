package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/telephrag/errlist"
)

func getLogFile(path string) *os.File {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		_, err := os.Create(path)
		if err != nil {
			log.Fatal(errlist.New(err))
		}
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Fatal(errlist.New(err))
		}
		return f
	}
	return f
}

func main() {

	logFile := getLogFile("/data/log.log")
	defer logFile.Close()
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetOutput(logFile)

	log.Print(errlist.New(nil).Set("event", "application_startup"))

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
	time.Sleep(time.Second * 60)
	log.Print(errlist.New(nil).Set("event", "application_graceful_shutdown"))

	// `docker stop %container%` has a timeout that can be set with `--time` flag.
	// Docker shall wait until timeout and than container will cease to exist.
	// Program must perform a cleanup within the said timeout to safely exit.
}
