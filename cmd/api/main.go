package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Server is running at :7979")
	http.ListenAndServe(":7979", nil)
}
