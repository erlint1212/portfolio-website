package views

import (
	_ "embed" // Blank import to enable embedding
)

// Using pragma (Compiler Directive)

//go:embed css/output.css
var cssContent string

func LoadCSS() string {
	return cssContent
}
