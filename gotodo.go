package main

import (
	"errors"
	"flag"
	"os"
)

type todoItem struct {
	item	string
}

var items = []todoItem{}

func manageTodoItems() ([]todoItem, error) {
	// validate that correct number of arguments is being received
	if len(os.Args) < 1 {
		return []todoItem{}, errors.New("Insufficient number of arguments")
	}

	action := flag.String("action", "view", "Action to perform on to-do list")

	flag.Parse()

	if !(*action == "add" || *action == "view" || *action == "complete") {
		return []todoItem{}, errors.New("Only 'view', 'add', and 'complete' actions are supported")
	}

	if (*action == "add") {
		item := flag.Arg(0)
		return []todoItem{todoItem{item}}, nil
	} else if (*action == "complete") {
		item := flag.Arg(0)
		return []todoItem{}, nil
	}
	return []todoItem{}, nil
}

func main() {

}