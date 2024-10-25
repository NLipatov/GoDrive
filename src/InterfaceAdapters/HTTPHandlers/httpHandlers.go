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

func (h *HttpHandler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, err := w.Write([]byte("OK"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
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

	const maxSizeBytes int64 = 300 * 1024 * 1024

	if err := r.ParseMultipartForm(maxSizeBytes); err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	uploadedFiles := r.MultipartForm.File["payload"]
	if len(uploadedFiles) != 1 {
		http.Error(w, "Send 1 file in a form-data with a 'payload' key", http.StatusBadRequest)
		return
	}

	formFile := uploadedFiles[0]
	f, err := formFile.Open()
	if err != nil {
		http.Error(w, "Error retrieving formFile from form data", http.StatusBadRequest)
		return
	}
	defer f.Close()

	fileData, err := io.ReadAll(f)
	if err != nil {
		http.Error(w, "Failed to read file bytes", http.StatusInternalServerError)
		return
	}

	file := Domain.File{
		Id:        uuid.Nil,
		Data:      fileData,
		ExpiresAt: time.Now().UTC().Add(15 * time.Hour),
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
