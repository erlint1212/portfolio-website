package views

import (
	"os"
	"log"
)

// No error handling because it will be embeded directly into html
func UnsafeLoadCSS() string {
	data, err := os.ReadFile("assets/css/output.css")
	if err != nil {
		log.Printf("[WARNING] Could not load CSS: %v", err)
		return "" 
	}

	return string(data)
}
