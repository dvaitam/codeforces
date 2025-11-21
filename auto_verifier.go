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
	"strings"
	"sync"
)

type verifierTask struct {
	problemPath  string
	solutionPath string
	verifierPath string
	logPath      string
	label        string
}

func main() {
	var concurrency int
	flag.IntVar(&concurrency, "c", 5, "concurrency level")
	flag.Parse()

	if concurrency < 1 {
		concurrency = 1
	}

	logDir := "~/log/verifier"
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		fmt.Printf("Error creating log directory: %v\n", err)
		return
	}

	var tasks []verifierTask
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		name := d.Name()
		if !strings.HasPrefix(name, "problem") || !strings.HasSuffix(name, ".txt") {
			return nil
		}
		index := strings.TrimPrefix(strings.TrimSuffix(name, ".txt"), "problem")
		if index == "" {
			return nil
		}
		dir := filepath.Dir(path)
		contest := filepath.Base(dir)
		solutionPath := filepath.Join(dir, contest+index+".go")
		if _, err := os.Stat(solutionPath); err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			fmt.Printf("Error checking %s: %v\n", solutionPath, err)
			return nil
		}

		verifierPath := filepath.Join(dir, "verifier"+index+".go")
		if _, err := os.Stat(verifierPath); err == nil {
			return nil
		} else if !os.IsNotExist(err) {
			fmt.Printf("Error checking %s: %v\n", verifierPath, err)
			return nil
		}

		logFile := filepath.Join(logDir, contest+index+"_verifier.log")
		tasks = append(tasks, verifierTask{
			problemPath:  path,
			solutionPath: solutionPath,
			verifierPath: verifierPath,
			logPath:      logFile,
			label:        contest + index,
		})
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No verifier tasks to process.")
		return
	}

	taskCh := make(chan verifierTask)
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			runVerifierWorker(workerID, taskCh)
		}(i + 1)
	}

	for _, t := range tasks {
		taskCh <- t
	}
	close(taskCh)

	wg.Wait()
	fmt.Println("All verifier tasks completed.")
}

var verifierSessionPattern = regexp.MustCompile(`session id:\s*([a-f0-9-]+)`)

func runVerifierWorker(workerID int, taskCh <-chan verifierTask) {
	var sessionID string
	for t := range taskCh {
		fmt.Printf("[verifier worker %d] Starting codex for %s\n", workerID, t.label)
		cmdStr := fmt.Sprintf("Write a go verifier for %s and save it as %s, use %s ref solution for testing , if the problem has mulitple solutions it should verify any possible solution do not rely on static input and ouput if multiple answers are correct, some times the problem might say print it any order, for testing float outputs make sure tolerance is according to the problem description, clean up any built artifacts after task finished", t.problemPath, t.verifierPath, t.solutionPath)
		args := buildVerifierCodexArgs(sessionID, cmdStr)
		cmd := exec.Command("codex", args...)
		output, err := cmd.CombinedOutput()
		if writeErr := os.WriteFile(t.logPath, output, 0o644); writeErr != nil {
			fmt.Printf("Error writing log for %s: %v\n", t.problemPath, writeErr)
		}
		if err != nil {
			fmt.Printf("Error running command for %s: %v\nLog saved to: %s\n", t.problemPath, err, t.logPath)
			continue
		}
		fmt.Printf("Successfully generated verifier %s\nLog saved to: %s\n", t.verifierPath, t.logPath)
		if newID := extractVerifierSessionID(output); newID != "" && newID != sessionID {
			sessionID = newID
			fmt.Printf("[verifier worker %d] Captured Codex session %s\n", workerID, sessionID)
		}
	}
}

func buildVerifierCodexArgs(sessionID, prompt string) []string {
	args := []string{"exec"}
	if sessionID == "" {
		args = append(args, "--full-auto")
	} else {
		args = append(args, "resume", sessionID)
	}
	args = append(args, prompt)
	return args
}

func extractVerifierSessionID(output []byte) string {
	matches := verifierSessionPattern.FindSubmatch(bytes.ToLower(output))
	if len(matches) < 2 {
		return ""
	}
	return string(matches[1])
}
