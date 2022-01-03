package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

func addTodoItem(db *bolt.DB, item string, date time.Time) error {
	itemBytes, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("could not marshal entry json: %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte("DB")).Bucket([]byte("TODOENTRIES")).Put([]byte(date.Format(time.RFC3339)), []byte(itemBytes))
		if err != nil {
			return fmt.Errorf("could not insert entry: %v", err)
		}
		return nil
	})
	return err
}

func completeTodoItem(db *bolt.DB, item string, date time.Time) error {
	var foundKey string
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("DB")).Bucket([]byte("TODOENTRIES"))
		bucket.ForEach(func(key, value []byte) error {
			if string(value)[1:len(string(value))-1] == item {
				foundKey = string(key)
			}
			return nil
		})
		return nil
	})
	if foundKey == "" {
		return nil
	}
	err = db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte("DB")).Bucket([]byte("TODOENTRIES")).Delete([]byte(foundKey))
		if err != nil {
			return fmt.Errorf("could not complete entry: %v", err)
		}
		return nil
	})
	return err
}

func displayTodoItems(db *bolt.DB) error {
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("DB")).Bucket([]byte("TODOENTRIES"))
		var i int = 1
		bucket.ForEach(func(key, value []byte) error {
			fmt.Printf("%d.\t%s\n", i, string(value))
			i++
			return nil
		})
		return nil
	})
	return err
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

	var err error

	if *action == "add" {
		item := flag.Arg(0)
		err = addTodoItem(db, item, time.Now())
	} else if *action == "complete" {
		item := flag.Arg(0)
		err = completeTodoItem(db, item, time.Now())
	}
	err = displayTodoItems(db)
	return err
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
