// Package main is the entrypoint for the todo application.
package main

import (
	"fmt"

	"github.com/user/todo/internal/model"
)

func main() {
	fmt.Println("Todo app")
	_ = model.PriorityMedium
	_ = model.StatusPending
}
