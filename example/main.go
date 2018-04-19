package main

import (
	"io"
	"net/http"
	"os"

	"github.com/phogolabs/parcello"
	_ "github.com/phogolabs/parcello/example/public"
)

func main() {
	file, err := parcello.Open("message.txt")
	if err != nil {
		panic(err)
	}

	if _, err = io.Copy(os.Stdout, file); err != nil {
		panic(err)
	}

	http.ListenAndServe(":8080", http.FileServer(parcello.Root("/website")))
}