package main

import (
	"flag"
	"os"
	"testing"
)

func Test_manageTodoCommands(t *testing.T) {
	tests := []struct {
		name    string   // name of test
		wantErr bool     // whether or not an error should be thrown
		osArgs  []string // command arguments used for test
	}{
		{"No parameters", false, []string{"cmd"}},
		{"Default parameters", false, []string{"cmd"}},
		{"Invalid action", true, []string{"cmd", "-action=foo", "item"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualOsArgs := os.Args
			defer func() {
				os.Args = actualOsArgs
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			}()

			os.Args = tt.osArgs
			err := manageTodoCommands()
			if (err != nil) != tt.wantErr {
				t.Errorf("getFileData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
