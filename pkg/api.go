package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

// ** API HANDLER ** //
func main() {
	code := `
def hello_world():
    print("Hello, world!")

hello_world()
`
	encodedCode := url.QueryEscape(code)
	data := url.Values{}
	data.Set("code", encodedCode)

	resp, err := http.Post("http://localhost:8080/execute", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Print the response from the server
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}

func executeHandler(w http.ResponseWriter, r *http.Request) {
	language := r.URL.Query().Get("language")

	if language == "" {
		http.Error(w, "Missing 'language' parameter", http.StatusBadRequest)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	encodedCode := r.FormValue("code")
	decodedCode, err := url.QueryUnescape(encodedCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileName := "code.py"
	err = os.WriteFile(fileName, []byte(decodedCode), 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the code
	cmd := exec.Command("python3", fileName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(output)
}
