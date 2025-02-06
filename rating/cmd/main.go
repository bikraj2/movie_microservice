package main

import (
	"fmt"
	"net/http"

	"bikraj.movie_microservice.net/rating/internal/controller/rating"
	httphandler "bikraj.movie_microservice.net/rating/internal/handler/http"
	"bikraj.movie_microservice.net/rating/internal/repository/memory"
)

func main() {
  fmt.Printf("Listening on port: 8082", a ...any)
	rep := memory.New()
	ctrl := rating.New(rep)
	h := httphandler.New(ctrl)

	http.Handle("/rating", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}

}
