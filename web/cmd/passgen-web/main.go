package main

import (
	"fmt"
	"net/http"
	"os"

	"codeberg.org/jhot/passgen/web"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("passgen-web listening on :%s\n", port)
	if err := http.ListenAndServe(":"+port, web.NewServer()); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
