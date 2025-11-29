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

type embedTask struct {
	testcasesPath string
	verifierPath  string
	solutionPath  string
	label         string
	logPath       string
	solBin        string
	verBin        string
}

func main() {
	var concurrency int
	flag.IntVar(&concurrency, "c", 3, "concurrency / codex session count")
	flag.Parse()

	if concurrency < 1 {
		concurrency = 1
	}

	logDir := filepath.Join("/home/ubuntu", "log", "embed_verifier")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		fmt.Printf("Error creating log directory: %v\n", err)
		return
	}

	var tasks []embedTask
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		name := d.Name()
		if !strings.HasPrefix(name, "testcases") || !strings.HasSuffix(name, ".txt") {
			return nil
		}

		dir := filepath.Dir(path)
		base := filepath.Base(dir)
		suffix := strings.TrimPrefix(strings.TrimSuffix(name, ".txt"), "testcases")

		verifierPath := filepath.Join(dir, "verifier"+suffix+".go")
		if _, err := os.Stat(verifierPath); err != nil {
			if !os.IsNotExist(err) {
				fmt.Printf("Error checking verifier %s: %v\n", verifierPath, err)
			}
			return nil
		}

		solutionPath := filepath.Join(dir, base+suffix+".go")
		if _, err := os.Stat(solutionPath); err != nil {
			if !os.IsNotExist(err) {
				fmt.Printf("Error checking solution %s: %v\n", solutionPath, err)
			}
			return nil
		}

		label := base + suffix
		logPath := filepath.Join(logDir, label+".log")
		tasks = append(tasks, embedTask{
			testcasesPath: path,
			verifierPath:  verifierPath,
			solutionPath:  solutionPath,
			label:         label,
			logPath:       logPath,
			solBin:        base + suffix,
			verBin:        base + "verifier" + suffix,
		})
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No testcases*.txt files with matching verifier/solution found.")
		return
	}

	sort.Slice(tasks, func(i, j int) bool { return tasks[i].label < tasks[j].label })

	taskCh := make(chan embedTask)
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			runEmbedWorker(workerID, taskCh)
		}(i + 1)
	}

	for _, t := range tasks {
		taskCh <- t
	}
	close(taskCh)

	wg.Wait()
	fmt.Println("All embed tasks completed.")
}

var sessionPattern = regexp.MustCompile(`session id:\s*([a-f0-9-]+)`)

func runEmbedWorker(workerID int, taskCh <-chan embedTask) {
	var sessionID string
	for t := range taskCh {
		fmt.Printf("[worker %d] Starting codex for %s\n", workerID, t.label)
		cmdStr := buildPrompt(t)
		args := buildCodexArgs(sessionID, cmdStr)
		cmd := exec.Command("codex", args...)
		output, err := cmd.CombinedOutput()
		if writeErr := os.WriteFile(t.logPath, output, 0o644); writeErr != nil {
			fmt.Printf("Error writing log for %s: %v\n", t.label, writeErr)
		}
		if err != nil {
			fmt.Printf("Error running codex for %s: %v\nLog saved to: %s\n", t.label, err, t.logPath)
		} else {
			fmt.Printf("Successfully processed %s\nLog saved to: %s\n", t.label, t.logPath)
		}
		if newID := extractSessionID(output); newID != "" && newID != sessionID {
			sessionID = newID
			fmt.Printf("[worker %d] Captured Codex session %s\n", workerID, sessionID)
		}
	}
}

func buildPrompt(t embedTask) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Verifier %s depends on %s, remove that dependency and embed those testcases as string data in the verifier itself. ", t.verifierPath, t.testcasesPath)
	fmt.Fprintf(&b, "It should also not depend on oracle or ref solution; embed everything it needs. ")
	fmt.Fprintf(&b, "Embed the code of %s directly into %s and refactor if needed. ", t.solutionPath, t.verifierPath)
	fmt.Fprintf(&b, "Try building the ref solution at /tmp/%s like `go build -o /tmp/%s %s`. ", t.solBin, t.solBin, t.solutionPath)
	fmt.Fprintf(&b, "Also build the verifier `go build -o /tmp/%s %s` and make sure `/tmp/%s /tmp/%s` returns with zero exit code. ", t.verBin, t.verifierPath, t.verBin, t.solBin)
	fmt.Fprintf(&b, "After the run, delete any /tmp artifacts like /tmp/%s and /tmp/%s (and related go-build* dirs) so /tmp stays clean. ", t.verBin, t.solBin)
	fmt.Fprintf(&b, "Once `/tmp/%s /tmp/%s` succeeds, delete the testcase file %s to remove the external dependency.", t.verBin, t.solBin, t.testcasesPath)
	return b.String()
}

func buildCodexArgs(sessionID, prompt string) []string {
	args := []string{"exec"}
	if sessionID == "" {
		args = append(args, "--full-auto")
	} else {
		args = append(args, "resume", sessionID)
	}
	args = append(args, prompt)
	return args
}

func extractSessionID(output []byte) string {
	matches := sessionPattern.FindSubmatch(bytes.ToLower(output))
	if len(matches) < 2 {
		return ""
	}
	return string(matches[1])
}
