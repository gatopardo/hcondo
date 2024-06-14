package controller

import (
	"net/http"

//	"github.com/gatopardo/hcondo/app/shared/view"
	"hcondo/app/shared/view"
)

// AboutGET displays the About page
func AboutGET(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "about/about"
	v.Render(w)
}
