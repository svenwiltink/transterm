package main

import (
	"github.com/svenwiltink/transterm/app"
)


// Show a navigable tree view of the current directory.
func main() {

	application := app.NewApplication(app.Config{
		AccountName:    "swiltink",
		PrivateKeyPath: "transip.key",
	})

	application.Run()
}
