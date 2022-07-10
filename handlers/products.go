package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/mayhendrap/go-microservices/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p.getProducts(w, r)
		return
	case http.MethodPost:
		p.addProduct(w, r)
		return
	case http.MethodPut:
		regexp := regexp.MustCompile(`/([0-9]+)`)
		g := regexp.FindAllStringSubmatch(r.URL.Path, -1)

		p.l.Println("GOT REGEX PATH :", g[0][1])

		if len(g) != 1 {
			p.l.Println("Invalid URI more than one id")
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.l.Println("Invalid URI more than one capture group")
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(g[0][1])
		if err != nil {
			p.l.Println("Invalid URI unable to convert to number", g[0][1])
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		err = p.updateProduct(id, w, r)
		if err == data.ErrProductNotFound {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, "Product not found", http.StatusInternalServerError)
			return
		}
		return
	default:
		// catch all
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")
	lp := data.GetProducts()
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	prod := &data.Product{}
	err := prod.FromJson(r.Body)
	if err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
	}

	data.AddProduct(prod)
}

func (p Products) updateProduct(id int, w http.ResponseWriter, r *http.Request) error {
	p.l.Println("Handle PUT Product")

	prod := &data.Product{}
	err := prod.FromJson(r.Body)
	if err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
	}

	return data.UpdateProduct(id, prod)
}
