package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
)

// fixTask represents a verifier build failure pulled from run.log.
type fixTask struct {
	label        string // e.g. 2000-2999/.../2172/problemA.txt
	dir          string // directory containing verifier/problem files
	suffix       string // problem letter, e.g. A, C2
	verifierPath string
	logPath      string
	errBody      string
}

func main() {
	var (
		runLog     string
		concurrent int
		memMB      int
	)
	flag.StringVar(&runLog, "log", "run.log", "path to run.log with verifier build failures")
	flag.IntVar(&concurrent, "c", runtime.NumCPU(), "codex concurrency")
	flag.IntVar(&memMB, "memlimit-mb", 0, "memory limit per codex process (0 = unlimited)")
	flag.Parse()

	if concurrent < 1 {
		concurrent = 1
	}
	if memMB < 0 {
		memMB = 0
	}

	content, err := os.ReadFile(runLog)
	if err != nil {
		fmt.Printf("failed to read %s: %v\n", runLog, err)
		return
	}

	logDir := filepath.Join(os.Getenv("HOME"), "log", "fix_verifier_build")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		fmt.Printf("failed to create log dir: %v\n", err)
		return
	}

	tasks := collectTasks(string(content), logDir)
	if len(tasks) == 0 {
		fmt.Println("no verifier build failures found in log")
		return
	}

	sort.Slice(tasks, func(i, j int) bool { return tasks[i].label < tasks[j].label })

	taskCh := make(chan fixTask)
	var wg sync.WaitGroup
	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			runWorker(worker, taskCh, memMB)
		}(i + 1)
	}

	for _, t := range tasks {
		taskCh <- t
	}
	close(taskCh)
	wg.Wait()
	fmt.Println("all fix tasks dispatched")
}

// Regex captures the fail line and its indented build output.
var failRe = regexp.MustCompile(`(?m)^\[FAIL\]\s+(\S+problem([A-Z0-9]+)\.txt)\s+\(build verifier\):[^\n]*\n((?:    .*\n)*)`)

func collectTasks(log string, logDir string) []fixTask {
	seen := make(map[string]bool)
	var tasks []fixTask

	matches := failRe.FindAllStringSubmatch(log, -1)
	for _, m := range matches {
		label := m[1]
		suffix := m[2]
		body := strings.ReplaceAll(m[3], "\r\n", "\n")
		body = strings.TrimRight(body, "\n")
		body = strings.ReplaceAll(body, "\n    ", "\n")
		body = strings.TrimSpace(body)

		dir := filepath.Dir(label)
		verifierPath := filepath.Join(dir, "verifier"+suffix+".go")
		if _, err := os.Stat(verifierPath); err != nil {
			// Skip if verifier is absent; nothing to fix automatically.
			continue
		}
		if seen[verifierPath] {
			continue
		}
		seen[verifierPath] = true
		logPath := filepath.Join(logDir, strings.ReplaceAll(dir, string(filepath.Separator), "_")+"_"+suffix+".log")
		tasks = append(tasks, fixTask{
			label:        label,
			dir:          dir,
			suffix:       suffix,
			verifierPath: verifierPath,
			logPath:      logPath,
			errBody:      body,
		})
	}
	return tasks
}

func runWorker(workerID int, taskCh <-chan fixTask, memLimitMB int) {
	var sessionID string
	for t := range taskCh {
		fmt.Printf("[worker %d] fixing %s\n", workerID, t.label)
		prompt := buildPrompt(t)
		args := buildCodexArgs(sessionID, prompt)
		cmd := buildCodexCmd(args, memLimitMB)
		out, err := cmd.CombinedOutput()
		if writeErr := os.WriteFile(t.logPath, out, 0o644); writeErr != nil {
			fmt.Printf("[worker %d] failed to write log for %s: %v\n", workerID, t.label, writeErr)
		}
		if err != nil {
			fmt.Printf("[worker %d] codex error for %s: %v (log %s)\n", workerID, t.label, err, t.logPath)
		} else {
			fmt.Printf("[worker %d] ok: %s (log %s)\n", workerID, t.label, t.logPath)
		}
		if newID := extractSessionID(out); newID != "" && newID != sessionID {
			sessionID = newID
			fmt.Printf("[worker %d] captured codex session %s\n", workerID, sessionID)
		}
	}
}

func buildPrompt(t fixTask) string {
	var b strings.Builder
	fmt.Fprintf(&b, "You are in %s. Verifier build fails for %s.\n", t.dir, t.verifierPath)
	fmt.Fprintf(&b, "Build error:\n%s\n", t.errBody)
	fmt.Fprintf(&b, "Fix only the verifier (do not touch statements or database). Remove unused imports, ensure package is main, and fix bad imports that point to paths like problem directories causing \"is not in std\" by using local code or copying needed helpers instead of imports. Keep the verifier behavior the same.\n")
	fmt.Fprintf(&b, "If the verifier shells out to build an oracle/solution, make sure it builds from the current directory (relative paths like ./oracleX.go) so it does not land under GOPATH and trigger \"package ... is not in std\" errors. Inline needed code instead of importing paths outside this directory.\n")
	fmt.Fprintf(&b, "Make `go build %s` succeed when run inside %s. After edits run gofmt. Delete any stray verifier binaries (files named like verifier without extensions) when done.", filepath.Base(t.verifierPath), t.dir)
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

func buildCodexCmd(args []string, memLimitMB int) *exec.Cmd {
	if memLimitMB <= 0 {
		return exec.Command("codex", args...)
	}
	memKB := memLimitMB * 1024
	script := fmt.Sprintf("ulimit -v %d; exec codex \"$@\"", memKB)
	bashArgs := []string{"-c", script, "--"}
	bashArgs = append(bashArgs, args...)
	return exec.Command("bash", bashArgs...)
}

// The codex CLI prints "session id: <uuid>" when a new session is created.
var sessionPattern = regexp.MustCompile(`session id:\s*([a-f0-9-]+)`)

func extractSessionID(output []byte) string {
	m := sessionPattern.FindSubmatch(bytes.ToLower(output))
	if len(m) < 2 {
		return ""
	}
	return string(m[1])
}

// removeVerifierBins deletes stray verifier binaries like verifierC/verifierD without extensions.
func removeVerifierBins(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, "verifier") && !strings.Contains(name, ".") {
			path := filepath.Join(dir, name)
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}
	return nil
}
