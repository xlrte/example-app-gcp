package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/database", dbHandler)
	r.HandleFunc("/bucket", bucketHandler)
	r.HandleFunc("/publish", publishHandler)
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func publishHandler(writer http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	client, err := pubsub.NewClient(ctx, os.Getenv("GCP_PROJECT_ID"))
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error()))
		return
	}
	topic := client.Topic(fmt.Sprintf("the_topic-%s", os.Getenv("XLRTE_ENV")))
	topic.Publish(ctx, &pubsub.Message{
		Data: []byte("published from web"),
	})

}

func bucketHandler(writer http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	client, err := storage.NewClient(ctx)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error()))
		return
	}
	bucket := client.Bucket(fmt.Sprintf("baz-%s", os.Getenv("XLRTE_ENV")))

	if req.Method == "POST" {
		obj := bucket.Object("data")
		w := obj.NewWriter(ctx)
		if _, err := fmt.Fprintf(w, "This object contains text.\n"); err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(err.Error()))
			return
		}
		// Close, just like writing a file.
		if err := w.Close(); err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(err.Error()))
			return
		}
		return
	}

	if req.Method == "GET" {
		obj := bucket.Object("data")
		r, err := obj.NewReader(ctx)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(err.Error()))
			return
		}
		defer r.Close()
		if _, err := io.Copy(writer, r); err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(err.Error()))
			return
		}
		writer.WriteHeader(200)

	}

}

func rootHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(err.Error()))
			return
		}
		fmt.Println("Received via POST " + string(body))
		writer.WriteHeader(200)
		writer.Write([]byte("alles gut!"))
		return
	}
	if req.Method == "PUT" {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(err.Error()))
			return
		}
		fmt.Println("Received via PUT " + string(body))
		writer.WriteHeader(200)
		writer.Write([]byte("alles gut!"))
		return
	}
	writer.WriteHeader(200)
	writer.Write([]byte("alles gut!"))
}

func dbHandler(writer http.ResponseWriter, req *http.Request) {
	db, err := getDb()
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error()))
		return
	}
	err = db.Ping()
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error() + " - failed to ping db"))
		return
	}

	writer.WriteHeader(200)
	writer.Write([]byte("alles gut!"))
}

func getDb() (*sqlx.DB, error) {
	dbString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_my-pg-db_HOST"),
		5432,
		os.Getenv("DB_my-pg-db_USER"),
		os.Getenv("DB_my-pg-db_PASSWORD"),
		fmt.Sprintf("my-pg-db-%s", os.Getenv("XLRTE_ENV")),
	)
	db, err := sqlx.Open("postgres", dbString)
	if err != nil {
		fmt.Println("ERR: " + dbString)
		return nil, err
	}

	return db, nil
}
