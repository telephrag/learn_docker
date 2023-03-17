package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/telephrag/errlist"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

type Counter struct {
	C int `bson:"counter"`
}

type Handler struct {
	mc *mongo.Client
}

func (h *Handler) Handle(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "chirtkem shuzhi")
	s, ok := syscall.Getenv("SECRET")
	if !ok {
		fmt.Fprintln(rw, errlist.New(fmt.Errorf("%d", http.StatusInternalServerError)))
	}
	fmt.Fprintf(rw, "мон тынад секретъёсыд тодӥсько: %s\n", s)

	col := h.mc.Database("local").Collection("example")
	res := col.FindOneAndUpdate(context.Background(),
		bson.M{},
		bson.M{
			"$inc": bson.M{
				"counter": 1,
			},
		},
	)
	if res.Err() != nil {
		log.Panic(res.Err())
	}

	var c Counter
	if err := res.Decode(&c); err != nil {
		log.Panic(err)
	}

	fmt.Fprintf(rw, "бам пол ветлӥськиз: %d", c.C)
}

func main() {

	logFile := getLogFile("log.log")
	defer logFile.Close()
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetOutput(logFile)

	clientOptions := options.Client().ApplyURI(
		"mongodb://mongo:27017",
	)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect(context.TODO())
	if err = client.Ping(context.TODO(), nil); err != nil {
		log.Panic(err)
	}

	if err := client.Database("local").Collection("example").Drop(context.Background()); err != nil {
		log.Panic(err)
	}
	col := client.Database("local").Collection("example")
	if _, err := col.InsertOne(context.Background(), Counter{0}); err != nil {
		log.Panic(err)
	}

	log.Print(errlist.New(nil).Set("event", "application_startup"))

	h := Handler{client}
	http.HandleFunc("/", h.Handle)

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
	// time.Sleep(time.Second * 60)
	log.Print(errlist.New(nil).Set("event", "application_graceful_shutdown"))

	// `docker stop %container%` has a timeout that can be set with `--time` flag.
	// Docker shall wait until timeout and than container will cease to exist.
	// Program must perform a cleanup within the said timeout to safely exit.
}
