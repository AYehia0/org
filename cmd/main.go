/*
The main entry point should call the app package to start the application.
*/
package main

import (
	"log"

	"github.com/ayehia0/org/pkg"
)

func main() {
	err := pkg.StartApp()
	if err != nil {
		log.Fatalf("could not start the app: %v", err)
	}
}
