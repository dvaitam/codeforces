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

func collectTasks(log string, logDir string) []fixTask {
	seen := make(map[string]bool)
	var tasks []fixTask

	lines := strings.Split(log, "\n")
	for _, line := range lines {
		if !strings.Contains(line, "is not in std") || !strings.Contains(line, "package") {
			continue
		}
		path := extractPathFromStdLine(line)
		if path == "" {
			continue
		}
		dir := filepath.Dir(path)
		base := strings.TrimSuffix(filepath.Base(path), ".go")
		suffix := problemSuffix(base)
		if suffix == "" {
			continue
		}
		verifierPath := filepath.Join(dir, "verifier"+suffix+".go")
		if _, err := os.Stat(verifierPath); err != nil {
			continue
		}
		if seen[verifierPath] {
			continue
		}
		seen[verifierPath] = true
		logPath := filepath.Join(logDir, strings.ReplaceAll(dir, string(filepath.Separator), "_")+"_"+suffix+".log")
		tasks = append(tasks, fixTask{
			label:        filepath.Join(dir, "problem"+suffix+".txt"),
			dir:          dir,
			suffix:       suffix,
			verifierPath: verifierPath,
			logPath:      logPath,
			errBody:      strings.TrimSpace(line),
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
	fmt.Fprintf(&b, "Fix only the verifier (do not touch statements or database). Ensure package is main, and fix bad imports/paths that cause \"is not in std\" by making the verifier self-contained in this directory. Copy needed helper code locally instead of importing sibling problem paths.\n")
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

// extractPathFromStdLine pulls the source path from a line containing "package ... is not in std".
func extractPathFromStdLine(line string) string {
	const marker = "package "
	pos := strings.Index(line, marker)
	if pos < 0 {
		return ""
	}
	pos += len(marker)
	end := strings.Index(line[pos:], " is not in std")
	if end < 0 {
		return ""
	}
	path := strings.TrimSpace(line[pos : pos+end])
	path = strings.TrimPrefix(path, "/usr/local/go/src/")
	return path
}

// problemSuffix returns the portion of the file name after the numeric contest ID.
func problemSuffix(base string) string {
	i := 0
	for i < len(base) && base[i] >= '0' && base[i] <= '9' {
		i++
	}
	if i == 0 || i == len(base) {
		return ""
	}
	return base[i:]
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
