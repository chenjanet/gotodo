package main

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func Test_manageTodoItems(t *testing.T) {
	tests := []struct {
		name	string // name of test
		want	[]todoItem // todoItem instance that function should return
		wantErr	bool // whether or not an error should be thrown
		osArgs	[]string // command arguments used for test	
	} {
		{"No parameters", []todoItem{}, false, []string{"cmd"}},
		{"Default parameters", []todoItem{}, false, []string{"cmd"}},
		{"Invalid action", []todoItem{}, true, []string{"cmd", "-action=foo", "item"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualOsArgs := os.Args
			defer func() {
				os.Args = actualOsArgs
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			}()

			os.Args = tt.osArgs
			got, err := manageTodoItems()
			if (err != nil) != tt.wantErr {
				t.Errorf("getFileData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {  // Asserting whether or not we get the correct wanted value
				t.Errorf("getFileData() = %v, want %v", got, tt.want)
			}
		})
	}
}