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
	newItem := todoItem{item}
	items = append(items, newItem)
	for i := 0; i < len(items); i++ {
		fmt.Printf("%d.	%s\n", i+1, items[i].item)
	}
}

func completeTodoItem(item string) {
	for i := 0; i < len(items); i++ {
		if item == items[i].item {
			items = append(items[:i], items[i+1:]...)
			break
		}
	}
}

func displayTodoItems() {
	for i := 0; i < len(items); i++ {
		fmt.Printf("%d.	%s\n", i+1, items[i].item)
	}
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

	if (*action == "add" || *action == "complete") && len(os.Args) < 3 {
		return errors.New("'add' and 'complete' commands require a specified item")
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

func exitGracefully(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func main() {
	// display usage info when user enters --help option
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <item to add or complete>\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}

	// processing user command
	err := manageTodoCommands()

	if err != nil {
		exitGracefully(err)
	}
}
