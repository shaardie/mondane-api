package main

import (
	"fmt"
	"os"

	"github.com/shaardie/mondane-api/api"
)

func mainWithError() error {
	return api.Run()
}

func main() {
	if err := mainWithError(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
