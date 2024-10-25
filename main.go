package main

import (
	"GoDrive/src/Application"
	"GoDrive/src/Infrastructure"
	"GoDrive/src/InterfaceAdapters/HTTPHandlers"
	"net/http"
)

func main() {
	Infrastructure.NewRepositoryDaemon().WatchStorageDirectory()
	repository := &Infrastructure.FileRepository{}
	fileService := Application.NewFileService(repository)
	httpHandlers := HTTPHandlers.NewHttpHandler(fileService)

	http.HandleFunc("/health", httpHandlers.Health)
	http.HandleFunc("/get", httpHandlers.Get)
	http.HandleFunc("/save", httpHandlers.Save)
	http.HandleFunc("/delete", httpHandlers.Delete)

	err := http.ListenAndServe(":10001", nil)
	if err != nil {
		panic(err)
	}
}
