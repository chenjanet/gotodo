package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

type todoItem struct {
	item string
}

type todoItemEntry struct {
	item string `json:"item"`
}

var items = []todoItem{}

func addTodoItem(db *bolt.DB, item string) error {
	newItem := todoItem{item}
	items = append(items, newItem)
	for i := 0; i < len(items); i++ {
		fmt.Printf("%d.	%s\n", i+1, items[i].item)
	}
	itemBytes, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("could not marshal entry json: %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte("DB")).Bucket([]byte("TODOENTRIES")).Put([]byte("items"), []byte(itemBytes))
		if err != nil {
			return fmt.Errorf("could not insert entry: %v", err)
		}
		return nil
	})
	return err
}

func completeTodoItem(db *bolt.DB, item string) error {
	for i := 0; i < len(items); i++ {
		if item == items[i].item {
			items = append(items[:i], items[i+1:]...)
			break
		}
	}
	itemBytes, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("could not marshal entry json: %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte("DB")).Bucket([]byte("TODOENTRIES")).Put([]byte("items"), []byte(itemBytes))
		if err != nil {
			return fmt.Errorf("could not insert entry: %v", err)
		}
		return nil
	})
	return err
}

func displayTodoItems(db *bolt.DB) {
	for i := 0; i < len(items); i++ {
		fmt.Printf("%d.	%s\n", i+1, items[i].item)
	}
}

func manageTodoCommands(db *bolt.DB) error {
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
		err := addTodoItem(db, item)
		if err != nil {
			return err
		}
		return nil
	} else if *action == "complete" {
		item := flag.Arg(0)
		err := completeTodoItem(db, item)
		if err != nil {
			return err
		}
		return nil
	}
	displayTodoItems(db)
	return nil
}

func exitGracefully(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func setupDB() (*bolt.DB, error) {
	db, dbErr := bolt.Open("gotodo.db", 0600, nil)

	if dbErr != nil {
		return nil, fmt.Errorf("could not open db, %v", dbErr)
	}

	dbErr = db.Update(func(tx *bolt.Tx) error {
		root, bucketErr := tx.CreateBucketIfNotExists([]byte("DB"))
		if bucketErr != nil {
			return fmt.Errorf("could not create root bucket: %v", bucketErr)
		}
		_, bucketErr = root.CreateBucketIfNotExists([]byte("TODOENTRIES"))
		if bucketErr != nil {
			return fmt.Errorf("could not create todo entry bucket: %v", bucketErr)
		}
		return nil
	})
	if dbErr != nil {
		return nil, fmt.Errorf("could not set up buckets, %v", dbErr)
	}
	fmt.Println("DB setup complete")
	return db, nil
}

func main() {
	db, dbErr := setupDB()

	if dbErr != nil {
		exitGracefully(dbErr)
	}

	defer db.Close()

	// display usage info when user enters --help option
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <item to add or complete>\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}

	// processing user command
	err := manageTodoCommands(db)

	if err != nil {
		exitGracefully(err)
	}
}
