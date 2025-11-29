package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
)

type solTask struct {
	verifierPath string
	solutionPath string
	label        string
	logPath      string
	solBin       string
	verBin       string
}

func main() {
	var concurrency int
	flag.IntVar(&concurrency, "c", 3, "concurrency / codex session count")
	flag.Parse()

	if concurrency < 1 {
		concurrency = 1
	}

	logDir := filepath.Join("/home/ubuntu", "log", "embed_sol")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		fmt.Printf("Error creating log directory: %v\n", err)
		return
	}

	var tasks []solTask
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		name := d.Name()
		if !strings.HasPrefix(name, "verifier") || !strings.HasSuffix(name, ".go") {
			return nil
		}

		dir := filepath.Dir(path)
		base := filepath.Base(dir)
		suffix := strings.TrimPrefix(strings.TrimSuffix(name, ".go"), "verifier")

		solutionPath := filepath.Join(dir, base+suffix+".go")
		if _, err := os.Stat(solutionPath); err != nil {
			if !os.IsNotExist(err) {
				fmt.Printf("Error checking solution %s: %v\n", solutionPath, err)
			}
			return nil
		}

		label := base + suffix
		logPath := filepath.Join(logDir, label+".log")
		tasks = append(tasks, solTask{
			verifierPath: path,
			solutionPath: solutionPath,
			label:        label,
			logPath:      logPath,
			solBin:       base + suffix,
			verBin:       base + "verifier" + suffix,
		})
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No verifier/solution pairs found.")
		return
	}

	sort.Slice(tasks, func(i, j int) bool { return tasks[i].label < tasks[j].label })

	taskCh := make(chan solTask)
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			runSolWorker(workerID, taskCh)
		}(i + 1)
	}

	for _, t := range tasks {
		taskCh <- t
	}
	close(taskCh)

	wg.Wait()
	fmt.Println("All embed_sol tasks completed.")
}

var solSessionPattern = regexp.MustCompile(`session id:\s*([a-f0-9-]+)`)

func runSolWorker(workerID int, taskCh <-chan solTask) {
	var sessionID string
	for t := range taskCh {
		fmt.Printf("[sol worker %d] Checking %s\n", workerID, t.label)
		logBuf := &bytes.Buffer{}

		ok := buildAndRun(t, logBuf)
		if ok {
			fmt.Printf("[sol worker %d] %s passed, cleaned /tmp binaries\n", workerID, t.label)
			continue
		}

		fmt.Printf("[sol worker %d] Starting codex for %s\n", workerID, t.label)
		prompt := buildSolPrompt(t)
		args := buildSolCodexArgs(sessionID, prompt)
		cmd := exec.Command("codex", args...)
		out, err := cmd.CombinedOutput()
		logBuf.WriteString("\n--- codex output ---\n")
		logBuf.Write(out)

		if writeErr := os.WriteFile(t.logPath, logBuf.Bytes(), 0o644); writeErr != nil {
			fmt.Printf("Error writing log for %s: %v\n", t.label, writeErr)
		}

		if err != nil {
			fmt.Printf("Codex failed for %s: %v\nLog saved to: %s\n", t.label, err, t.logPath)
		} else {
			fmt.Printf("Codex succeeded for %s\nLog saved to: %s\n", t.label, t.logPath)
		}

		if newID := extractSolSessionID(out); newID != "" && newID != sessionID {
			sessionID = newID
			fmt.Printf("[sol worker %d] Captured Codex session %s\n", workerID, sessionID)
		}
	}
}

func buildAndRun(t solTask, logBuf *bytes.Buffer) bool {
	solBin := filepath.Join("/tmp", t.solBin)
	verBin := filepath.Join("/tmp", t.verBin)

	cleanup := func() {
		_ = os.Remove(solBin)
		_ = os.Remove(verBin)
	}
	defer cleanup()

	logBuf.WriteString(fmt.Sprintf("Building solution: go build -o %s %s\n", solBin, t.solutionPath))
	if out, err := exec.Command("go", "build", "-o", solBin, t.solutionPath).CombinedOutput(); err != nil {
		logBuf.Write(out)
		logBuf.WriteString(fmt.Sprintf("solution build failed: %v\n", err))
		return false
	} else {
		logBuf.Write(out)
	}

	logBuf.WriteString(fmt.Sprintf("Building verifier: go build -o %s %s\n", verBin, t.verifierPath))
	if out, err := exec.Command("go", "build", "-o", verBin, t.verifierPath).CombinedOutput(); err != nil {
		logBuf.Write(out)
		logBuf.WriteString(fmt.Sprintf("verifier build failed: %v\n", err))
		return false
	} else {
		logBuf.Write(out)
	}

	logBuf.WriteString(fmt.Sprintf("Running: %s %s\n", verBin, solBin))
	cmd := exec.Command(verBin, solBin)
	out, err := cmd.CombinedOutput()
	logBuf.Write(out)
	if err != nil {
		logBuf.WriteString(fmt.Sprintf("run failed: %v\n", err))
		return false
	}

	logBuf.WriteString("run succeeded\n")
	return true
}

func buildSolPrompt(t solTask) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Fix %s and %s so that the verifier run succeeds. ", t.verifierPath, t.solutionPath)
	fmt.Fprintf(&b, "It should not depend on an oracle or external ref solution; embed everything it needs. ")
	fmt.Fprintf(&b, "Embed the code of %s directly into %s and refactor if needed. ", t.solutionPath, t.verifierPath)
	fmt.Fprintf(&b, "Try building the ref solution at /tmp/%s like `go build -o /tmp/%s %s`. ", t.solBin, t.solBin, t.solutionPath)
	fmt.Fprintf(&b, "Also build the verifier `go build -o /tmp/%s %s` and make sure `/tmp/%s /tmp/%s` returns with zero exit code. ", t.verBin, t.verifierPath, t.verBin, t.solBin)
	fmt.Fprintf(&b, "After the run, delete any /tmp artifacts like /tmp/%s and /tmp/%s (and related go-build* dirs) so /tmp stays clean.", t.verBin, t.solBin)
	return b.String()
}

func buildSolCodexArgs(sessionID, prompt string) []string {
	args := []string{"exec"}
	if sessionID == "" {
		args = append(args, "--full-auto")
	} else {
		args = append(args, "resume", sessionID)
	}
	args = append(args, prompt)
	return args
}

func extractSolSessionID(output []byte) string {
	matches := solSessionPattern.FindSubmatch(bytes.ToLower(output))
	if len(matches) < 2 {
		return ""
	}
	return string(matches[1])
}
