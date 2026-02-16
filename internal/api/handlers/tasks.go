package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kamil-budzik/csv-processor/internal/db"
	"github.com/kamil-budzik/csv-processor/internal/models"
)

type Handler struct {
	repo  TaskRepository
	store FileStorage
}

func NewHandler(repo TaskRepository, store FileStorage) *Handler {
	if repo == nil {
		panic("NewHandler: TaskRepository cannot be nil")
	}
	if store == nil {
		panic("NewHandler: store cannnot be nil")
	}
	return &Handler{repo: repo, store: store}
}

func parseUUID(c *gin.Context, param string) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param(param))
	if err != nil {
		log.Printf("Failed to Parse param to UUID. Param: %s", param)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   "Invalid task_id format",
		})
		return uuid.UUID{}, false
	}
	return id, true
}

func (h *Handler) GetTask(c *gin.Context) {
	id, ok := parseUUID(c, "task_id")
	if !ok {
		return
	}

	ctx := c.Request.Context()
	task, err := h.repo.GetTask(ctx, id)

	if errors.Is(err, db.ErrTaskNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"data":   "Task not found",
		})
		return
	}
	if err != nil {
		log.Printf("DB Error inside GetTask %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"data":   "Internal Server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   task,
	})
}

func (h *Handler) GetAllTasks(c *gin.Context) {
	ctx := c.Request.Context()
	tasks, err := h.repo.GetTasks(ctx)
	if err != nil {

		log.Printf("DB Error inside GetAllTasks %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"data":   "Internal Server Error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   tasks,
	})
}

func (h *Handler) PostTask(c *gin.Context) {
	ctx := c.Request.Context()

	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Error in parsing formfile %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"data":   "Internal Server Error. Failed to form file",
		})
		return
	}

	opened, err := file.Open()
	id := uuid.New()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"data":   "Internal Server Error. Failed to open file",
		})
		return
	}

	defer opened.Close()

	filePath, err := h.store.UploadCSV(ctx, fmt.Sprintf("%s.csv", id.String()), file.Size, opened)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": "Failed to upload file to Minio"})
		return
	}

	task, err := h.repo.CreateTask(ctx, filePath, id)
	if err != nil {
		log.Printf("DB Error inside PostTask %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"data":   "Internal Server Error. Failed to create Task",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   task,
	})

}

func (h *Handler) PutTask(c *gin.Context) {
	id, ok := parseUUID(c, "task_id")
	if !ok {
		return
	}

	var input models.TaskUpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "data": err})
		return
	}

	ctx := c.Request.Context()
	task, err := h.repo.UpdateTask(ctx, id, input)
	if errors.Is(err, db.ErrTaskNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"data":   "Task not found",
		})
		return
	}

	if errors.Is(err, db.ErrInvalidTransition) {
		c.JSON(http.StatusConflict, gin.H{
			"status": "error",
			"data":   "Invalid Status Transition",
		})
		return
	}

	if err != nil {
		log.Printf("DB Error inside PutTask %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"data":   "Internal Server Error. Failed to Update Task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   task,
	})
}

func (h *Handler) DeleteTask(c *gin.Context) {
	id, ok := parseUUID(c, "task_id")
	if !ok {
		return
	}

	ctx := c.Request.Context()
	err := h.repo.DeleteTask(ctx, id)

	if errors.Is(err, db.ErrTaskNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"data":   "Task not found",
		})
		return
	}

	if err != nil {
		log.Printf("DB Error inside Delete Task %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "data": "Internal Server Error"})
		return
	}

	c.Status(http.StatusNoContent)
}
