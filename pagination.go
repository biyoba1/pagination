package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

type Pagination struct {
	TotalItems   int
	ItemsPerPage int
	Page         int
}

type PageData struct {
	Items      []string
	Pagination *Pagination
}

func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.ItemsPerPage
}

func (p *Pagination) GetLimit() int {
	return p.ItemsPerPage
}

func (p *Pagination) GetTotalPages() int {
	return int(math.Ceil(float64(p.TotalItems) / float64(p.ItemsPerPage)))
}

func subtract(x, y int) int {
	return x - y
}

func add(x, y int) int {
	return x + y
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			page = 1
		}

		pagination := &Pagination{
			TotalItems:   1000,
			ItemsPerPage: 1,
			Page:         page,
		}

		items := make([]string, 0, pagination.ItemsPerPage)
		for i := pagination.GetOffset(); i < pagination.GetOffset()+pagination.ItemsPerPage && i < pagination.TotalItems; i++ {
			items = append(items, fmt.Sprintf("Item %d", i+1))
		}
		
		log.Printf("ОТКРЫТА СТРАНИЦА - %d", page)
		
		pageData := &PageData{
			Items:      items,
			Pagination: pagination,
		}

		funcMap := template.FuncMap{
			"sub": subtract,
			"add": add,
		}

		t, err := template.New("page").Funcs(funcMap).Parse(`
			<html>
			<body>
			<h1>Page {{.Pagination.Page}} of {{.Pagination.GetTotalPages}}</h1>
			<ul>
			{{range .Items}}
				<li>{{.}}</li>
			{{end}}
			</ul>
			{{if gt .Pagination.Page 1}}
				<a href="?page={{.Pagination.Page | sub 1}}">Previous</a>
			{{end}}
			{{if lt .Pagination.Page .Pagination.GetTotalPages}}
				<a href="?page={{.Pagination.Page | add 1}}">Next</a>
			{{end}}
			</body>
			</html>
		`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, pageData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.ListenAndServe(":8080", nil)
}
