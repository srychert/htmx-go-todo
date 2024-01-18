package router

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"todo/model"
)

func Setup() {
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
		writeTodoRow(w)
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
		todoId, err := getTodoId(w, r)

		if err != nil {
			return
		}

		var doneUint uint8 = 0
		if done {
			doneUint = 1
		}

		err = model.SetTodoDone(todoId, doneUint)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	h5 := func(w http.ResponseWriter, r *http.Request) {
		todoColumn := r.Header.Get("X-Todo-Column")
		todoId, err := getTodoId(w, r)

		if err != nil {
			return
		}

		todo, err := model.GetTodoById(todoId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles(
			"templates/partials/todo-input.html",
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		todoInput := model.TodoInput{TodoId: todo.Id, Name: todoColumn, Value: todo.Title}

		err = tmpl.ExecuteTemplate(w, "todo-input", todoInput)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	h6 := func(w http.ResponseWriter, r *http.Request) {
		todoColumn := r.Header.Get("X-Todo-Column")
		value := r.FormValue(todoColumn)
		todoId, err := getTodoId(w, r)

		if err != nil {
			return
		}

		todo, err := model.GetTodoById(todoId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if todoColumn == "title" {
			todo.Title = value
			err = model.UpdateTodo(todo)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		writeTodoRow(w, todo)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", h1)
	http.HandleFunc("/add-todo/", h2)
	http.HandleFunc("/home/", h3)
	http.HandleFunc("/mark-todo/", h4)
	http.HandleFunc("/edit-todo/", h5)
	http.HandleFunc("/update-todo/", h6)

	log.Fatal(http.ListenAndServe(":8000", nil))
}

func getTodoId(w http.ResponseWriter, r *http.Request) (uint64, error) {
	todoIdString := r.Header.Get("X-Todo-Id")
	todoId, err := strconv.Atoi(todoIdString)

	if err != nil {
		http.Error(w, "Invalid todoId", http.StatusBadRequest)
		return 0, err
	}

	return uint64(todoId), err
}

func writeTodoRow(w http.ResponseWriter, todos ...model.Todo) {
	var todo model.Todo

	if len(todos) == 0 {
		todo, _ = model.CreateTodo("")
	} else {
		todo = todos[0]
	}

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
