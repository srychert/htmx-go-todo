package main

import (
	"html/template"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"strconv"
	"todo/model"
)

func main() {
	model.Setup()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	h1 := func(w http.ResponseWriter, r *http.Request) {
		todosData, _ := model.GetAllTodos()

		todos := map[string][]model.Todo{
			"Todos": todosData,
		}

		// layout file must be the first parameter in ParseFiles!
		templates, err := template.ParseFiles(
			"templates/layout.html",
			"templates/home.html",
			"templates/aside.html",
			"templates/todo-table.html",
			"templates/partials/todo-row.html",
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = templates.Execute(w, todos)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		todo, _ := model.CreateTodo("")

		tmpl, err := template.ParseFiles(
			"templates/partials/todo-row.html",
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "todo-row", todo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	h3 := func(w http.ResponseWriter, r *http.Request) {
		todosData, _ := model.GetAllTodos()

		todos := map[string][]model.Todo{
			"Todos": todosData,
		}

		templates, err := template.ParseFiles(
			"templates/home.html",
			"templates/todo-table.html",
			"templates/partials/todo-row.html",
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = templates.ExecuteTemplate(w, "home", todos)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	h4 := func(w http.ResponseWriter, r *http.Request) {
		done := r.FormValue("done") == "on"
		todoIdString := r.Header.Get("X-Todo-Id")

		todoId, err := strconv.Atoi(todoIdString)
		if err != nil {
			http.Error(w, "Invalid todoId", http.StatusBadRequest)
			return
		}

		var doneUint uint8 = 0
		if done {
			doneUint = 1
		}

		err = model.SetTodoDone(uint64(todoId), doneUint)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/add-todo/", h2)
	http.HandleFunc("/home", h3)
	http.HandleFunc("/update-todo", h4)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
