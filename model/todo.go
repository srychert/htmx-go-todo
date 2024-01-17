package model

type Todo struct {
	Id    uint64 `json:"id"`
	Title string `json:"title"`
	Date  string `json:"date"`
	Done  uint8  `json:"done"`
}

func CreateTodo(title string) (Todo, error) {
	statement := `insert into todo(title, date, done) values($1,  DATETIME('now'), $2) returning id;`

	id := 0
	err := db.QueryRow(statement, title, 0).Scan(&id)

	if err != nil {
		return Todo{}, err
	}

	todo, err := GetTodoById(uint64(id))
	if err != nil {
		return Todo{}, err
	}

	return todo, err
}

func GetAllTodos() ([]Todo, error) {
	var todos []Todo

	statement := `select id, title, date, done from todo;`

	rows, err := db.Query(statement)
	if err != nil {
		return todos, err
	}

	defer rows.Close()

	for rows.Next() {
		var id uint64
		var title string
		var date string
		var done uint8

		err := rows.Scan(&id, &title, &date, &done)
		if err != nil {
			return todos, err
		}
		todo := Todo{
			Id:    id,
			Title: title,
			Date:  date,
			Done:  done,
		}

		todos = append(todos, todo)
	}

	return todos, err
}

func GetTodoById(id uint64) (Todo, error) {
	todo := Todo{}
	todo.Id = id

	statement := `select title, date, done from todo where id = $1;`

	row, err := db.Query(statement, id)
	if err != nil {
		return todo, err
	}

	for row.Next() {
		var title string
		var date string
		var done uint8
		err := row.Scan(&title, &date, &done)
		if err != nil {
			return todo, err
		}

		todo.Title = title
		todo.Date = date
		todo.Done = done
	}
	return todo, err
}

func SetTodoDone(id uint64, done uint8) error {
	statement := `update todo set done=$2 where id=$1;`
	_, err := db.Query(statement, id, done)

	return err
}
