package main

import (
    "fmt"
    "log"
    "encoding/json"
    "net/http"
    "goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

var(
  db = "store"
  collection = "products"
)




type Product struct {
    Id      string `json:"id"`
    Name    string   `json:"name"`
    Title   string   `json:"title"`
    Description string `json:"description"`
    Price   string   `json:"price"`
}

func JSONErrorHandler(w http.ResponseWriter, message string, code int) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(code)
    fmt.Fprintf(w, "{message: %q}", message)
}

func JSONResponseHandler(w http.ResponseWriter, json []byte, code int) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(code)
    w.Write(json)
}

func main() {
    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
    }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)

    mux := goji.NewMux()
    mux.HandleFunc(pat.Get("/products"), getAllProducts(session))
    mux.HandleFunc(pat.Post("/product"), createProduct(session))
    mux.HandleFunc(pat.Get("/product/:id"), productById(session))
  //  mux.HandleFunc(pat.Put("/books/:isbn"), updateBook(session))
    mux.HandleFunc(pat.Delete("/product/:id"), deleteProduct(session))
    http.ListenAndServe("localhost:8080", mux)
}

func getAllProducts(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        c := session.DB(db).C(collection)

        var products []Product
        err := c.Find(bson.M{}).All(&products)
        if err != nil {
            JSONErrorHandler(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed get all books: ", err)
            return
        }

        respBody, err := json.MarshalIndent(products, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        JSONResponseHandler(w, respBody, http.StatusOK)
    }
}

func createProduct(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()
  log.Println("request ", r.Body)
        var product Product
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&product)
        if err != nil{
              log.Println("Failed create product : ", err)
            JSONErrorHandler(w, "Bad Request", http.StatusBadRequest)
            return
        }
        c := session.DB(db).C(collection)
        err = c.Insert(product)
        if err != nil {
            if mgo.IsDup(err) {
                JSONErrorHandler(w, "Product is already exists", http.StatusBadRequest)
                return
            }
            JSONErrorHandler(w, "Error in Database", http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
    }
}

func productById(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()
        productId := pat.Param(r, "id")
        c := session.DB(db).C(collection)
        var product Product
        err := c.Find(bson.M{"id": productId}).One(&product)
        if err != nil {
            JSONErrorHandler(w, "Error in Database", http.StatusInternalServerError)
            return
        }

        if product.Id == "" {
            JSONErrorHandler(w, "Product Id doesn't exist", http.StatusNotFound)
            return
        }

        respBody, err := json.MarshalIndent(product, "", "  ")
        if err != nil {
            log.Fatal(err)
        }
        JSONResponseHandler(w, respBody, http.StatusOK)
    }
}

func deleteProduct(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        productId := pat.Param(r, "id")

        c := session.DB(db).C(collection)

        err := c.Remove(bson.M{"id": productId})
        if err != nil {
            switch err {
            default:
                JSONErrorHandler(w, "Error in Database", http.StatusInternalServerError)
                return
            case mgo.ErrNotFound:
                JSONErrorHandler(w, "Product Id doesn't exist", http.StatusNotFound)
                return
            }
        }

        w.WriteHeader(http.StatusNoContent)
    }
}


/*
func updateBook(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        isbn := pat.Param(r, "isbn")

        var book Book
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&book)
        if err != nil {
            ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
            return
        }

        c := session.DB("store").C("books")

        err = c.Update(bson.M{"isbn": isbn}, &book)
        if err != nil {
            switch err {
            default:
                ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
                log.Println("Failed update book: ", err)
                return
            case mgo.ErrNotFound:
                ErrorWithJSON(w, "Book not found", http.StatusNotFound)
                return
            }
        }

        w.WriteHeader(http.StatusNoContent)
    }
}


*/
