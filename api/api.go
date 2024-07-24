package api

import (
    "encoding/json"
    "net/http"
    "runix/executors"
)

// SetupRoutes initializes the API routes
func SetupRoutes(r *mux.Router) {
    r.HandleFunc("/execute", executeHandler).Methods("POST")
}

func executeHandler(w http.ResponseWriter, r *http.Request) {
    var requestBody map[string]string
    if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

	code, ok := requestBody["code"]
	if !ok {
        http.Error(w, "Missing 'code' field", http.StatusBadRequest)
        return
    }

	language, ok := requestBody["language"]

	if !ok {
		http.Error(w, "Missing 'Language' field", http.StatusBadRequest)
	}
	if language == "python"{
		cmd, err := exec.Command("/bin/sh", "executors/python.sh")

		if err!=nil{
			fmt.Println("Error %s",err)

		}
	}

	// handle response
	response := map[string]string{
        "result": result,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)

