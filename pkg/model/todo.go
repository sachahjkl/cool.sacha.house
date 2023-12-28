package model

import (
	"fmt"

	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title  string `gorm:"not null"`
	Done   bool
	UserID uint
}

func AddTodo(db *gorm.DB, title string, done bool, userId uint) (*Todo, error) {

	if len(title) == 0 {
		return nil, fmt.Errorf("todo title can't be empty")
	}

	todo := Todo{
		Title:  title,
		Done:   done,
		UserID: userId,
	}

	// insert new todo
	err := db.Create(&todo).Error
	return &todo, err
}

func DeleteTodo(db *gorm.DB, id uint, userId uint) error {
	return db.Where("user_id = ?", userId).Delete(&Todo{}, id).Error
}

func ToggleTodo(db *gorm.DB, id uint, userId uint) (*Todo, error) {
	var todo Todo
	err := db.Where("user_id = ?", userId).First(&todo, id).Error
	if err != nil {
		return nil, err
	}

	// toggle the todo
	todo.Done = !todo.Done

	err = db.Save(&todo).Error
	return &todo, err
}

func AllTodos(db *gorm.DB, userId uint) (*[]Todo, error) {
	var todos []Todo
	err := db.Where("user_id = ?", userId).Order("created_at desc, title").Find(&todos).Error
	return &todos, err
}
