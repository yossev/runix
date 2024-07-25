package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
)

// ** API HANDLER ** //

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/execute", executeHandler).Methods("GET")

	// Start the server
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", router)
}

// handler called when excute endpoint is hit
func executeHandler(w http.ResponseWriter, r *http.Request) {
	language := r.URL.Query().Get("language")
	code := r.URL.Query().Get("code")

	if language == "" || code == "" {
		http.Error(w, "Missing 'language' or 'code' parameter", http.StatusBadRequest)
		return
	}
	// THIS FUNCTION ISNT APPLYING SOLID YET!, FUNCTION ALSO HANDLES FILE CREATION
	if language == "Python" {
		tmpFile, err := os.Create("file.py")
		if err != nil {
			fmt.Println("%s: Cannot create temp file", err)
		}
		defer tmpFile.Close()
		defer os.Remove(tmpFile.Name())

		// write the code in the temp file
		if _, err := tmpFile.Write([]byte(code)); err != nil {
			http.Error(w, "Could not write to temporary file", http.StatusInternalServerError)
			return
		}
		cmd := exec.Command("./executors/python.sh", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		if err != nil {
			http.Error(w, fmt.Sprintf("Execution error: %s", err), http.StatusInternalServerError)
			return
		}
		w.Write(output)
	}
}
