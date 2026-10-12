package db

import (
	"database/sql"
	"errors"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

const taskColumns = `id, date, title, comment, repeat`

func AddTask(task *Task) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetTask(id string) (*Task, error) {
	task := &Task{}
	row := db.QueryRow(`SELECT `+taskColumns+` FROM scheduler WHERE id = ?`, id)
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("задача не найдена")
	}
	if err != nil {
		return nil, err
	}
	return task, nil
}

func Tasks(limit int) ([]*Task, error) {
	return selectTasks(`SELECT `+taskColumns+` FROM scheduler ORDER BY date LIMIT ?`, limit)
}

func TasksBySearch(limit int, search string) ([]*Task, error) {
	like := "%" + search + "%"
	return selectTasks(
		`SELECT `+taskColumns+` FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`,
		like, like, limit,
	)
}

func TasksByDate(limit int, date string) ([]*Task, error) {
	return selectTasks(`SELECT `+taskColumns+` FROM scheduler WHERE date = ? LIMIT ?`, date, limit)
}

func selectTasks(query string, args ...any) ([]*Task, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*Task{}
	for rows.Next() {
		task := &Task{}
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func UpdateTask(task *Task) error {
	res, err := db.Exec(
		`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID,
	)
	if err != nil {
		return err
	}
	return checkRowsAffected(res)
}

func UpdateDate(id, date string) error {
	res, err := db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, date, id)
	if err != nil {
		return err
	}
	return checkRowsAffected(res)
}

func DeleteTask(id string) error {
	res, err := db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return checkRowsAffected(res)
}

func checkRowsAffected(res sql.Result) error {
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}
