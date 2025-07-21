package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type contestInfo struct {
	ID       string
	Path     string
	Problems []string
}

var contests map[string]*contestInfo

var indexTmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html><body>
<h1>Contests</h1>
<ul>
{{range .}}
<li><a href="/contest/{{.ID}}">{{.ID}}</a></li>
{{end}}
</ul>
</body></html>`))

var contestTmpl = template.Must(template.New("contest").Parse(`
<!DOCTYPE html>
<html><body>
<h1>Contest {{.ID}}</h1>
<ul>
{{range .Problems}}
<li><a href="/contest/{{$.ID}}/problem/{{.}}">Problem {{.}}</a></li>
{{end}}
</ul>
</body></html>`))

var problemTmpl = template.Must(template.New("problem").Parse(`
<!DOCTYPE html>
<html><body>
<h1>Contest {{.Contest}} Problem {{.Letter}}</h1>
<pre>{{.Statement}}</pre>
<form action="/contest/{{.Contest}}/problem/{{.Letter}}/submit" method="post" enctype="multipart/form-data">
<select name="lang">
<option value="c">C</option>
<option value="cpp">C++</option>
<option value="java">Java</option>
<option value="python">Python</option>
<option value="go">Go</option>
<option value="rust">Rust</option>
</select><br>
<textarea name="code" rows="20" cols="80"></textarea><br>
<input type="file" name="file"><br>
<input type="submit" value="Submit">
</form>
</body></html>`))

var resultTmpl = template.Must(template.New("result").Parse(`
<!DOCTYPE html>
<html><body>
<h1>Result for Contest {{.Contest}} Problem {{.Letter}}</h1>
<pre>{{.Output}}</pre>
<a href="/contest/{{.Contest}}/problem/{{.Letter}}">Back</a>
</body></html>`))

func scanContests(root string) (map[string]*contestInfo, error) {
	result := make(map[string]*contestInfo)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		base := filepath.Base(path)
		if _, err := strconv.Atoi(base); err == nil {
			entries, err := os.ReadDir(path)
			if err != nil {
				return nil
			}
			var probs []string
			for _, e := range entries {
				name := e.Name()
				if strings.HasPrefix(name, "problem") && strings.HasSuffix(name, ".txt") {
					letter := strings.TrimSuffix(strings.TrimPrefix(name, "problem"), ".txt")
					probs = append(probs, letter)
				}
			}
			if len(probs) > 0 {
				sort.Strings(probs)
				result[base] = &contestInfo{ID: base, Path: path, Problems: probs}
			}
		}
		return nil
	})
	return result, err
}

func findVerifier(dir, letter string) string {
	cand := filepath.Join(dir, "verifier"+letter+".go")
	if _, err := os.Stat(cand); err == nil {
		return cand
	}
	cand = filepath.Join(dir, "verifier.go")
	if _, err := os.Stat(cand); err == nil {
		return cand
	}
	return ""
}

func compileSource(srcPath, lang string) (string, string, error) {
	tmpDir, err := os.MkdirTemp("", "submit")
	if err != nil {
		return "", "", err
	}
	exe := filepath.Join(tmpDir, "main")
	var cmd *exec.Cmd
	switch lang {
	case "c":
		cmd = exec.Command("gcc", srcPath, "-O2", "-std=c11", "-o", exe)
	case "cpp":
		cmd = exec.Command("g++", srcPath, "-O2", "-std=c++17", "-o", exe)
	case "go":
		cmd = exec.Command("go", "build", "-o", exe, srcPath)
	case "rust":
		cmd = exec.Command("rustc", "-O", srcPath, "-o", exe)
	case "java":
		javaDir := filepath.Join(tmpDir, "java")
		if err := os.Mkdir(javaDir, 0755); err != nil {
			return "", "", err
		}
		cmd = exec.Command("javac", "-d", javaDir, srcPath)
		exe = filepath.Join(tmpDir, "run-java.sh")
		script := fmt.Sprintf("#!/bin/sh\njava -cp %s Main \"$@\"\n", javaDir)
		if err := os.WriteFile(exe, []byte(script), 0755); err != nil {
			return "", "", err
		}
	case "python":
		exe = filepath.Join(tmpDir, "run-python.sh")
		script := fmt.Sprintf("#!/bin/sh\npython3 %s \"$@\"\n", srcPath)
		if err := os.WriteFile(exe, []byte(script), 0755); err != nil {
			return "", "", err
		}
		return exe, "", nil
	default:
		return "", "", fmt.Errorf("unknown language")
	}

	if cmd != nil {
		out, err := cmd.CombinedOutput()
		if err != nil {
			return "", string(out), err
		}
		if lang == "java" {
			// ensure wrapper is executable
			return exe, string(out), nil
		}
	}
	return exe, "", nil
}

func submitSolution(w http.ResponseWriter, r *http.Request, c *contestInfo, letter string) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	lang := r.FormValue("lang")
	var data []byte
	file, _, err := r.FormFile("file")
	if err == nil {
		defer file.Close()
		data, _ = io.ReadAll(file)
	} else {
		data = []byte(r.FormValue("code"))
	}
	extMap := map[string]string{"c": ".c", "cpp": ".cpp", "java": ".java", "python": ".py", "go": ".go", "rust": ".rs"}
	ext := extMap[lang]
	if ext == "" {
		http.Error(w, "unknown language", http.StatusBadRequest)
		return
	}
	srcPath := filepath.Join(c.Path, "user"+strings.ToUpper(letter)+ext)
	if err := os.WriteFile(srcPath, data, 0644); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	exe, compileOut, err := compileSource(srcPath, lang)
	output := bytes.Buffer{}
	if err != nil {
		output.WriteString("Compilation failed:\n")
		output.WriteString(compileOut)
		output.WriteString(err.Error())
	} else {
		verifier := findVerifier(c.Path, letter)
		if verifier != "" {
			cmd := exec.Command("go", "run", verifier, exe)
			cmd.Dir = c.Path
			res, err := cmd.CombinedOutput()
			output.Write(res)
			if err != nil {
				output.WriteString("\nVerifier error: " + err.Error())
			}
		} else {
			output.WriteString("Compiled successfully. No verifier available.")
		}
	}
	resultTmpl.Execute(w, map[string]string{
		"Contest": c.ID,
		"Letter":  letter,
		"Output":  output.String(),
	})
}

func problemPage(w http.ResponseWriter, r *http.Request, c *contestInfo, letter string) {
	stmtPath := filepath.Join(c.Path, "problem"+letter+".txt")
	data, err := os.ReadFile(stmtPath)
	if err != nil {
		http.Error(w, "problem not found", http.StatusNotFound)
		return
	}
	problemTmpl.Execute(w, map[string]string{
		"Contest":   c.ID,
		"Letter":    letter,
		"Statement": string(data),
	})
}

func contestPage(w http.ResponseWriter, r *http.Request, cid string) {
	c := contests[cid]
	if c == nil {
		http.NotFound(w, r)
		return
	}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/contest/"+cid), "/")
	if len(parts) >= 3 && parts[1] == "problem" && parts[2] != "" {
		if len(parts) == 4 && parts[3] == "submit" {
			submitSolution(w, r, c, parts[2])
			return
		}
		if r.Method == http.MethodGet {
			problemPage(w, r, c, parts[2])
			return
		}
	}
	if r.URL.Path != "/contest/"+cid {
		http.NotFound(w, r)
		return
	}
	contestTmpl.Execute(w, c)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	ids := make([]string, 0, len(contests))
	for id := range contests {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	var list []*contestInfo
	for _, id := range ids {
		list = append(list, contests[id])
	}
	indexTmpl.Execute(w, list)
}

func main() {
	var err error
	contests, err = scanContests(".")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/contest/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/contest/"), "/")
		if len(parts) == 0 || parts[0] == "" {
			http.NotFound(w, r)
			return
		}
		contestPage(w, r, parts[0])
	})
	http.ListenAndServe(":8080", nil)
}
