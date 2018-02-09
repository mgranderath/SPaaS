package views

import (
	"fmt"
	"net/http"
)

// HomePage is the / page
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Placeholder")
}
