package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var formTmpl = template.Must(template.New("form").Parse(`
<!DOCTYPE html>
<html>
<head><title>Add Solution</title></head>
<body>
<h1>Add Solution for Contest {{.ContestID}}</h1>
<form action="/contest/{{.ContestID}}/add-solution" method="post" enctype="multipart/form-data">
<label>Problem Letter: <input type="text" name="problem" required></label><br>
<label>Language:
<select name="lang">
  <option value="CPP">CPP</option>
  <option value="Java">Java</option>
  <option value="Go">Go</option>
  <option value="Python">Python</option>
</select>
</label><br>
<label>Code:</label><br>
<textarea name="code" rows="20" cols="80"></textarea><br>
<label>Or Upload File: <input type="file" name="file"></label><br>
<input type="submit" value="Submit">
</form>
</body>
</html>
`))

func contestDir(cid string) (string, error) {
	n, err := strconv.Atoi(cid)
	if err != nil {
		return "", err
	}
	sn := n
	thousands := (n / 1000) * 1000
	tDir := fmt.Sprintf("%d-%d", thousands, thousands+999)
	n = n % 1000
	hundreds := (n / 100) * 100
	hDir := fmt.Sprintf("%d-%d", thousands+hundreds, thousands+hundreds+99)
	n = n % 100
	tens := (n / 10) * 10
	teDir := fmt.Sprintf("%d-%d", thousands+hundreds+tens, thousands+hundreds+tens+9)
	cDir := fmt.Sprintf("%d", sn)
	path := filepath.Join(tDir, hDir, teDir, cDir)
	return path, nil
}

func addSolution(w http.ResponseWriter, r *http.Request, cid string) {
	if r.Method == http.MethodGet {
		formTmpl.Execute(w, map[string]string{"ContestID": cid})
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	lang := r.FormValue("lang")
	letter := r.FormValue("problem")

	var data []byte
	file, _, err := r.FormFile("file")
	if err == nil {
		defer file.Close()
		data, _ = io.ReadAll(file)
	} else {
		data = []byte(r.FormValue("code"))
	}

	extMap := map[string]string{"CPP": ".cpp", "Java": ".java", "Go": ".go", "Python": ".py"}
	ext := extMap[lang]
	if ext == "" {
		http.Error(w, "unknown language", http.StatusBadRequest)
		return
	}

	path, err := contestDir(cid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = os.MkdirAll(path, 0755)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(path, "sol"+strings.ToUpper(letter)+ext)
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Saved to %s", filePath)
}

func handler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) == 3 && parts[0] == "contest" && parts[2] == "add-solution" {
		addSolution(w, r, parts[1])
		return
	}
	http.NotFound(w, r)
}

func main() {
	http.HandleFunc("/contest/", handler)
	http.ListenAndServe(":8080", nil)
}
