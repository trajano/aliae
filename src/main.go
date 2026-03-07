package main

import "github.com/jandedobbeleer/aliae/src/app"

var (
	Version = "development"
)

func main() {
	app.Execute(Version)
}
