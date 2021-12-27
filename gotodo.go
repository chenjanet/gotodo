package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type todoItem struct {
	item string
}

var items = []todoItem{}

func addTodoItem(item string) {
	fmt.Print(item)
}

func completeTodoItem(item string) {
	fmt.Print(item)
}

func displayTodoItems() {

}

func manageTodoCommands() error {
	// validate that correct number of arguments is being received
	if len(os.Args) < 1 {
		return errors.New("Insufficient number of arguments")
	}

	action := flag.String("action", "view", "Action to perform on to-do list")

	flag.Parse()

	if !(*action == "add" || *action == "view" || *action == "complete") {
		return errors.New("Only 'view', 'add', and 'complete' actions are supported")
	}

	if *action == "add" {
		item := flag.Arg(0)
		addTodoItem(item)
		return nil
	} else if *action == "complete" {
		item := flag.Arg(0)
		completeTodoItem(item)
		return nil
	}
	displayTodoItems()
	return nil
}

func main() {

}
