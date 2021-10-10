package router

import (
	"book-restapi/book"
	"book-restapi/configs"
	"book-restapi/customerror"
	"book-restapi/mongolib"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// for get the title from body
type Title struct {
	Name string `json:"title"`
}

/* globally decaring variables*/
var appConfigs *configs.App
var mongoConfigs *configs.Mongo
var ctx context.Context
var datastore *mongolib.MongoDatastore

func init() {
	// all glabal variables are initializing inside the init function
	appConfigs, mongoConfigs, err := configs.ReadConfigs()
	datastore, err = mongolib.NewDatasore(mongoConfigs)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println(appConfigs.Name)
	fmt.Println(datastore.Config.Collections)
	fmt.Println(datastore.Session)
	fmt.Println(ctx)

	if err != nil {
		log.Fatal()
	}
}

/*routers for api*/
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	books, err := mongolib.GetBooks(ctx, (datastore.DB.Client().Database("library").Collection("books")))

	if err != nil {
		json.NewEncoder(w).Encode(&customerror.CustomError{Error: err, StatusCode: 500})
		return
	}

	json.NewEncoder(w).Encode(&books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	book, err := mongolib.GetBook(ctx, (datastore.DB.Client().Database("library").Collection("books")),
		params["id"])

	if err != nil {
		json.NewEncoder(w).Encode(&customerror.CustomError{Error: err, StatusCode: 500})
		return
	}
	json.NewEncoder(w).Encode(&book)
}

func insertBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book book.Book

	err := json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		json.NewEncoder(w).Encode(&customerror.CustomError{Error: err, StatusCode: 500})
		return
	}

	res, err := mongolib.InsertBook(ctx, (datastore.DB.Client().Database("library").Collection("books")),
		book)
	if err != nil && res != nil {
		json.NewEncoder(w).Encode(&customerror.CustomError{Error: err, StatusCode: 500})
		return
	}

	json.NewEncoder(w).Encode(&book)

}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	//title := r.Body
	var title Title
	err := json.NewDecoder(r.Body).Decode(&title)

	if err != nil {
		json.NewEncoder(w).Encode(&customerror.CustomError{Error: err, StatusCode: 500})
		return
	}

	res, err := mongolib.UpdateBook(ctx, (datastore.DB.Client().Database("library").Collection("books")), params["id"], title.Name)

	if err != nil {
		json.NewEncoder(w).Encode(&customerror.CustomError{Error: err, StatusCode: 500})
		return
	}

	json.NewEncoder(w).Encode(&res)

}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var title Title

	err := json.NewDecoder(r.Body).Decode(&title)
	fmt.Println(title.Name)

	if err != nil {
		json.NewEncoder(w).Encode(&customerror.CustomError{Error: err, StatusCode: 500})
		return
	}

	res, err := mongolib.DeleteBook(ctx, (datastore.DB.Client().Database("library").Collection("books")), title.Name)

	json.NewEncoder(w).Encode(&res)
}

func RouterInitializer() {

	r := mux.NewRouter()

	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", insertBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books", deleteBook).Methods("DELETE")

	err := http.ListenAndServe(":8000", r)

	if err != nil {
		log.Fatal(err)
	}

}
