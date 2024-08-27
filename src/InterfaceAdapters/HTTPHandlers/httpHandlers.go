package HTTPHandlers

import (
	"GoDrive/src/Application"
	"GoDrive/src/Domain"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
)

type HttpHandler struct {
	fileService *Application.FileService
}

func NewHttpHandler(fileService *Application.FileService) *HttpHandler {
	return &HttpHandler{fileService: fileService}
}

func (h *HttpHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	file, err := h.fileService.Get(id)
	if err != nil {
		http.Error(w, "File not found", http.StatusBadRequest)
		return
	}

	_, err = w.Write(file.Data)
	if err != nil {
		http.Error(w, "Failed to write file data", http.StatusInternalServerError)
		return
	}
}

func (h *HttpHandler) Save(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	file := Domain.File{
		Id:        uuid.Nil,
		Data:      bytes,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 15),
	}

	id, err := h.fileService.Save(file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(id.String()))
}

func (h *HttpHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	_, err = h.fileService.Delete(id)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
