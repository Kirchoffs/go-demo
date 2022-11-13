package api

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/google/uuid"
    "github.com/gorilla/mux"
)

type Item struct {
    ID uuid.UUID `json:"id"`
    Name string `json:"name"`
}

type Server struct {
    *mux.Router
    shoppingItems []Item
}

func NewServer() *Server {
    s := &Server{
        Router: mux.NewRouter(),
        shoppingItems: []Item{},
    }
    s.routes()
    return s
}

func (s *Server) routes() {
    s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, world!"))
    }).Methods("GET")

    s.HandleFunc("/test", s.test()).Methods("GET")
    
    s.HandleFunc("/shopping-items", s.listShoppingItems()).Methods("GET")
    s.HandleFunc("/shopping-items", s.createShoppingItem()).Methods("POST")
    s.HandleFunc("/shopping-items/{id}", s.removeShoppingItem()).Methods("DELETE")
}

func (s *Server) test() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("Test Start")

        context := r.Context()

        select {
        case <-time.After(8 * time.Second):
            fmt.Println("Test End")
            w.Write([]byte("Test End"))
        
        case <-context.Done():
            fmt.Println("Connection Closed")
            err := context.Err()
            fmt.Println(err)
        }
    }
}

func (s *Server) createShoppingItem() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var i Item 
        if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        i.ID = uuid.New()
        s.shoppingItems = append(s.shoppingItems, i)

        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(i); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
    }
}

func (s *Server) listShoppingItems() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(s.shoppingItems); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
    }
}

func (s *Server) removeShoppingItem() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id, err := uuid.Parse(mux.Vars(r)["id"])
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return 
        }

        for i, item := range s.shoppingItems {
            if item.ID == id {
                s.shoppingItems = append(s.shoppingItems[:i], s.shoppingItems[i+1:]...)
                break
            }
        }

        w.WriteHeader(http.StatusNoContent)
    }
}
