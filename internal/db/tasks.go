package db

type Task struct {
	TaskID string `json:"task_id"`
}

func GetTasks() ([]Task, error) {
	rows, err := DB.Query("SELECT task_id FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.TaskID); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
