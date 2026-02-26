package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/lucidiusss/todo-fiber-react/database"
	"github.com/lucidiusss/todo-fiber-react/models"
)

func GetTasks(c fiber.Ctx) error {
	var tasks []models.Task

	result := database.DB.Find(&tasks)

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

	// remove fields that shouldn't be updated
	delete(updates, "id")
	delete(updates, "created_at")
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
		// check if task already exists
		var count int64

		database.DB.Model(&models.Task{}).Where("LOWER(title) = LOWER(?) AND id != ?", title, id).Count(&count)

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

	var task models.Task

	result := database.DB.First(&task, id)

	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Task not found",
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
