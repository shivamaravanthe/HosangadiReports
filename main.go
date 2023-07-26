package main

import (
	"log"
	"shivamaravanthe/HosangadiReports/constants"
	"shivamaravanthe/HosangadiReports/database"
	"shivamaravanthe/HosangadiReports/server"
)

func main() {
	database.ConnectDB()
	log.Printf("Server started on Port %s", constants.PORT)
	server.Server()
}
