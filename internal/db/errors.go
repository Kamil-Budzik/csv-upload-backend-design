package db

import "errors"

var ErrTaskNotFound = errors.New("Task not found")
var ErrInvalidTransition = errors.New("Invalid Task transition")
