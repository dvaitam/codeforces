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

type task struct {
	problemPath string
	goPath      string
	logPath     string
	label       string
}

func main() {
	var concurrency int
	flag.IntVar(&concurrency, "c", 5, "concurrency level")
	flag.Parse()

	if concurrency < 1 {
		concurrency = 1
	}

	logDir := "~/log"
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		fmt.Printf("Error creating log directory: %v\n", err)
		return
	}

	var tasks []task
	err = filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
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
		goFile := filepath.Join(dir, contest+index+".go")
		if _, statErr := os.Stat(goFile); statErr == nil {
			return nil
		} else if !os.IsNotExist(statErr) {
			fmt.Printf("Error checking %s: %v\n", goFile, statErr)
			return nil
		}
		logFile := filepath.Join(logDir, contest+index+".log")
		tasks = append(tasks, task{
			problemPath: path,
			goPath:      goFile,
			logPath:     logFile,
			label:       contest + index,
		})
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No matching files found to process.")
		return
	}

	taskCh := make(chan task)
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			runSolutionWorker(workerID, taskCh)
		}(i + 1)
	}

	for _, t := range tasks {
		taskCh <- t
	}
	close(taskCh)

	wg.Wait()
	fmt.Println("All tasks completed.")
}

var sessionPattern = regexp.MustCompile(`session id:\s*([a-f0-9-]+)`)

func runSolutionWorker(workerID int, taskCh <-chan task) {
	var sessionID string
	for t := range taskCh {
		fmt.Printf("[worker %d] Starting codex for %s\n", workerID, t.label)
		cmdStr := fmt.Sprintf("write a go solution for %s and save it as %s, clean up any built artifacts after task finished", t.problemPath, t.goPath)
		args := buildCodexArgs(sessionID, cmdStr)
		cmd := exec.Command("codex", args...)
		output, err := cmd.CombinedOutput()
		if writeErr := os.WriteFile(t.logPath, output, 0644); writeErr != nil {
			fmt.Printf("Error writing log for %s: %v\n", t.problemPath, writeErr)
		}
		if err != nil {
			fmt.Printf("Error running command for %s: %v\nLog saved to: %s\n", t.problemPath, err, t.logPath)
		} else {
			fmt.Printf("Successfully processed %s to %s\nLog saved to: %s\n", t.problemPath, t.goPath, t.logPath)
		}
		if newID := extractSessionID(output); newID != "" && newID != sessionID {
			sessionID = newID
			fmt.Printf("[worker %d] Captured Codex session %s\n", workerID, sessionID)
		}
	}
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
