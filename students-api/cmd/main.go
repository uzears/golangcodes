package main

import (
	"net/http"

	"github.com/uzears/students-api/internal/config"
)

func main() {

	cfg = config.MustLoad()

	router := http.NewServeMux()

	http.ListenAndServe(cfg.HTTPServer.Addr, router)

}
