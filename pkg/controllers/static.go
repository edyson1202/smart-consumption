package controllers

import (
	"net/http"
)

func StaticHandler(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, nil)
	}
}

func FAQ(tpl Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   string
	}{
		{
			Question: "How tall is mount Everest?",
			Answer:   "It is 8848 meters tall.",
		},
		{
			Question: "When did we go to the moon?",
			Answer:   "It 1969.",
		},
		{
			Question: "What rocket did we use?",
			Answer:   "A SaturnV.",
		},
		{
			Question: "Who was the first man to step on the moon?",
			Answer:   "It was Neil Armstrong.",
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, questions)
	}
}
