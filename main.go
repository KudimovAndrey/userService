package main

import (
	"net/http"
	"userService/service"
)

const addr string = "localhost:8080"

func main() {
	mux := http.NewServeMux()
	srv := service.NewService()
	mux.HandleFunc("/", srv.Handle)
	mux.HandleFunc("/makeFriends", srv.MakeFriends)
	http.ListenAndServe(addr, mux)
}
