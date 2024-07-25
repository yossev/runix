package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
)

// ** API HANDLER ** //

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/execute", executeHandler).Methods("POST")

	// Start the server
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", router)
}

// handler called when execute endpoint is hit
func executeHandler(w http.ResponseWriter, r *http.Request) {
	language := r.URL.Query().Get("language")

	if language == "" {
		http.Error(w, "Missing 'language' parameter", http.StatusBadRequest)
		return
	}

	// Parse multipart form to handle file upload
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("code")
	if err != nil {
		http.Error(w, "Unable to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	tmpFile, err := os.CreateTemp("", "code_*.py") // Create a temp file with appropriate extension
	if err != nil {
		http.Error(w, "Unable to create temp file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	// Write the uploaded code to the temp file
	_, err = io.Copy(tmpFile, file)
	if err != nil {
		http.Error(w, "Unable to write to temp file", http.StatusInternalServerError)
		return
	}

	// Ensure the file is readable
	if err := os.Chmod(tmpFile.Name(), 0600); err != nil {
		http.Error(w, "Unable to set file permissions", http.StatusInternalServerError)
		return
	}

	// Execute the code based on the language
	var cmd *exec.Cmd
	if language == "Python" {
		cmd = exec.Command("./executors/python.sh", tmpFile.Name())
	} else {
		http.Error(w, "Unsupported language", http.StatusBadRequest)
		return
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Execution error: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(output)
}
