package main

import (
	"flag"
	"fmt"
	"net/http"
	"userService/service"
)

func main() {
	var addr string
	flag.StringVar(&addr, "host", "localhost:8080", "server connection address")
	flag.Parse()
	mux := http.NewServeMux()
	srv := service.NewService()
	mux.HandleFunc("/", srv.Handle)
	mux.HandleFunc("/makeFriends", srv.MakeFriends)
	fmt.Printf("server started on %s\n", addr)
	http.ListenAndServe(addr, mux)
}
