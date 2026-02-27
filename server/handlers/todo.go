package handlers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/lucidiusss/todo-fiber-react/database"
	"github.com/lucidiusss/todo-fiber-react/models"
	"gorm.io/gorm"
)

func GetTasks(c fiber.Ctx) error {
	userID := c.Locals("user_id")

	fmt.Printf("userID: %v (type: %T)\n", userID, userID)

	var tasks []models.Task

	result := database.DB.Where("user_id = ?", userID).Find(&tasks)

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch tasks",
			"error":   result.Error.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Tasks fetched successfully",
		"data":    tasks,
		"count":   len(tasks),
	})
}

func CreateTask(c fiber.Ctx) error {

	task := new(models.Task)

	// parse req body
	if err := c.Bind().Body(task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// validate required fields
	if task.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Title is required",
		})
	}

	// check if task already exists
	var count int64

	database.DB.Model(&models.Task{}).Where("LOWER(title) = LOWER(?) AND id != ? AND user_id = ?", task.Title, task.ID, task.UserID).Count(&count)

	if count > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error":   "Task with this title already exists",
		})
	}

	// create task in db
	result := database.DB.Create(&task)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create task",
			"error":   result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Task created successfully",
		"data":    task,
	})
}

func GetTask(c fiber.Ctx) error {
	id := c.Params("id")

	var task models.Task

	result := database.DB.First(&task, id)

	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Task not found",
		})
	}

	return c.Status(fiber.StatusFound).JSON(fiber.Map{
		"success": true,
		"message": "Task by id " + id + " was found",
		"data":    task,
	})
}

func UpdateTask(c fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id")

	// parse updates into map
	var updates map[string]any
	if err := c.Bind().Body(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// check if task exists
	var task models.Task
	if err := database.DB.First(&task, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Task not found",
		})
	}

	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Task not found or you don't have permission to update it",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Database error",
			"error":   err.Error(),
		})
	}

	// Remove protected fields
	delete(updates, "id")
	delete(updates, "user_id")
	delete(updates, "created_at")
	delete(updates, "updated_at")
	delete(updates, "deleted_at")

	// validate field types
	if title, ok := updates["title"]; ok {
		if str, ok := title.(string); ok {
			if len(str) < 1 {
				return c.Status(400).JSON(fiber.Map{
					"success": false,
					"message": "Title cannot be empty",
				})
			}
			if len(str) > 255 {
				return c.Status(400).JSON(fiber.Map{
					"success": false,
					"message": "Title too long (max 255 chars)",
				})
			}
		} else {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Title must be a string",
			})
		}

		userID := c.Locals("user_id")

		// check if task already exists
		var count int64

		database.DB.Model(&models.Task{}).Where("LOWER(title) = LOWER(?) AND id != ? AND user_id = ?", title, id, userID).Count(&count)

		if count > 0 {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false,
				"message": "Task with this title already exists",
				"error":   "DUPLICATE_TITLE",
			})
		}

	}

	if completed, ok := updates["completed"]; ok {
		if _, ok := completed.(bool); !ok {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Completed must be a boolean",
			})
		}
	}

	// perform update - only fields in the map are updated

	if err := database.DB.Model(&task).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to update task",
			"error":   err.Error(),
		})
	}

	database.DB.First(&task, id)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Task updated successfully",
		"data":    task,
	})

}

func DeleteTask(c fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id")

	var task models.Task
	if err := database.DB.First(&task, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Task not found",
		})
	}

	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Task not found or you don't have permission to update it",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Database error",
			"error":   err.Error(),
		})
	}

	//delete todo
	database.DB.Delete(&task)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Todo deleted successfully",
		"data":    task,
	})
}
