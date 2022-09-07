package main

import (
	"github.com/TharinduBalasooriya/goRestApi/controller"
	"github.com/TharinduBalasooriya/goRestApi/logger"
	_ "github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	var router = http.NewServeMux()

	//Home Toute

	//Other Routes
	router.HandleFunc("/api/movies", controller.GetAllMovies)

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8081", logger.LogRequestHandler(router)))

}
