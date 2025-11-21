package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type solutionTask struct {
	key        string
	problemRel string
	goRel      string
	logPath    string
	onComplete func(string, bool)
}

type verifierTask struct {
	key         string
	problemRel  string
	solutionRel string
	verifierRel string
	logPath     string
}

type orchestrator struct {
	root            string
	watcher         *fsnotify.Watcher
	solutionPool    *workerPool[solutionTask]
	verifierPool    *workerPool[verifierTask]
	solutionLogDir  string
	verifierLogDir  string
	mu              sync.Mutex
	activeSolutions map[string]bool
	pendingVerifier map[string]verifierTask
}

func main() {
	var root string
	var maxSolutions int
	var maxVerifiers int

	flag.StringVar(&root, "root", ".", "root directory to watch")
	flag.IntVar(&maxSolutions, "solution-max", 5, "maximum concurrent solution workers")
	flag.IntVar(&maxVerifiers, "verifier-max", 3, "maximum concurrent verifier workers")
	flag.Parse()

	absRoot, err := filepath.Abs(root)
	if err != nil {
		log.Fatalf("failed to resolve root %s: %v", root, err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to resolve home directory: %v", err)
	}
	solutionLogDir := filepath.Join(homeDir, "log")
	verifierLogDir := filepath.Join(solutionLogDir, "verifier")
	if err := os.MkdirAll(solutionLogDir, 0o755); err != nil {
		log.Fatalf("failed to create log directory %s: %v", solutionLogDir, err)
	}
	if err := os.MkdirAll(verifierLogDir, 0o755); err != nil {
		log.Fatalf("failed to create verifier log directory %s: %v", verifierLogDir, err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("failed to create watcher: %v", err)
	}
	defer watcher.Close()

	orch := &orchestrator{
		root:            absRoot,
		watcher:         watcher,
		solutionLogDir:  solutionLogDir,
		verifierLogDir:  verifierLogDir,
		activeSolutions: make(map[string]bool),
		pendingVerifier: make(map[string]verifierTask),
	}

	orch.solutionPool = newWorkerPool("solution", maxSolutions, orch.solutionRunner)
	orch.verifierPool = newWorkerPool("verifier", maxVerifiers, orch.verifierRunner)

	if err := orch.addWatchRecursive(absRoot); err != nil {
		log.Fatalf("failed to watch directory %s: %v", absRoot, err)
	}

	log.Printf("Watching %s for new problems/solutions...", absRoot)

	done := make(chan struct{})
	go orch.eventLoop(done)

	<-done
}

func (o *orchestrator) eventLoop(done chan struct{}) {
	defer close(done)
	for {
		select {
		case event, ok := <-o.watcher.Events:
			if !ok {
				return
			}
			o.handleEvent(event)
		case err, ok := <-o.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("watcher error: %v", err)
		}
	}
}

func (o *orchestrator) addWatchRecursive(path string) error {
	return filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if err := o.watcher.Add(p); err != nil {
			return fmt.Errorf("add watch for %s: %w", p, err)
		}
		return nil
	})
}

func (o *orchestrator) handleEvent(event fsnotify.Event) {
	if event.Name == "" {
		return
	}
	if event.Op&(fsnotify.Create|fsnotify.Rename) == 0 {
		return
	}

	info, err := os.Stat(event.Name)
	if err == nil && info.IsDir() {
		if err := o.addWatchRecursive(event.Name); err != nil {
			log.Printf("failed adding watch for %s: %v", event.Name, err)
		}
		return
	}

	if contest, index, ok := parseProblemFile(event.Name); ok {
		o.scheduleSolution(event.Name, contest, index)
		return
	}

	if contest, index, ok := parseSolutionFile(event.Name); ok {
		o.scheduleVerifier(event.Name, contest, index)
	}
}

func (o *orchestrator) scheduleSolution(problemPath, contest, index string) {
	key := contest + index
	dir := filepath.Dir(problemPath)
	goPath := filepath.Join(dir, fmt.Sprintf("%s%s.go", contest, index))
	if fileExists(goPath) {
		return
	}

	relProblem := o.relative(problemPath)
	relGo := o.relative(goPath)
	logPath := filepath.Join(o.solutionLogDir, key+".log")

	o.mu.Lock()
	if o.activeSolutions[key] {
		o.mu.Unlock()
		return
	}
	o.activeSolutions[key] = true
	o.mu.Unlock()

	log.Printf("Queued solution task for %s (%s)", key, relProblem)
	o.solutionPool.Submit(solutionTask{
		key:        key,
		problemRel: relProblem,
		goRel:      relGo,
		logPath:    logPath,
		onComplete: o.onSolutionComplete,
	})
}

func (o *orchestrator) scheduleVerifier(solutionPath, contest, index string) {
	key := contest + index
	dir := filepath.Dir(solutionPath)
	problemPath := filepath.Join(dir, fmt.Sprintf("problem%s.txt", index))
	verifierPath := filepath.Join(dir, fmt.Sprintf("verifier%s.go", index))

	if !fileExists(problemPath) || fileExists(verifierPath) {
		return
	}

	task := verifierTask{
		key:         key,
		problemRel:  o.relative(problemPath),
		solutionRel: o.relative(solutionPath),
		verifierRel: o.relative(verifierPath),
		logPath:     filepath.Join(o.verifierLogDir, key+"_verifier.log"),
	}

	o.mu.Lock()
	if o.activeSolutions[key] {
		if _, exists := o.pendingVerifier[key]; !exists {
			o.pendingVerifier[key] = task
			log.Printf("Verifier for %s queued until solution completes", key)
		}
		o.mu.Unlock()
		return
	}
	o.mu.Unlock()

	log.Printf("Queued verifier task for %s", key)
	o.verifierPool.Submit(task)
}

func (o *orchestrator) onSolutionComplete(key string, success bool) {
	o.mu.Lock()
	delete(o.activeSolutions, key)
	pending, hasPending := o.pendingVerifier[key]
	if hasPending {
		delete(o.pendingVerifier, key)
	}
	o.mu.Unlock()

	if hasPending {
		if success {
			log.Printf("Dispatching queued verifier for %s", key)
			o.verifierPool.Submit(pending)
		} else {
			log.Printf("Dropping queued verifier for %s because solution failed", key)
		}
	}
}

func (o *orchestrator) solutionRunner(w *worker[solutionTask], task solutionTask) {
	prompt := fmt.Sprintf("write a go solution for %s and save it as %s, clean up any built artifacts after task finished", task.problemRel, task.goRel)
	success := o.runCodexCommand("solution", w.id, &w.sessionID, prompt, task.logPath)
	if task.onComplete != nil {
		task.onComplete(task.key, success)
	}
}

func (o *orchestrator) verifierRunner(w *worker[verifierTask], task verifierTask) {
	prompt := fmt.Sprintf("Write a go verifier for %s and save it as %s, use %s ref solution for testing , if the problem has mulitple solutions it should verify any possible solution do not rely on static input and ouput if multiple answers are correct, some times the problem might say print it any order, for testing float outputs make sure tolerance is according to the problem description, clean up any built artifacts after task finished", task.problemRel, task.verifierRel, task.solutionRel)
	o.runCodexCommand("verifier", w.id, &w.sessionID, prompt, task.logPath)
}

func (o *orchestrator) runCodexCommand(kind string, workerID int, session *string, prompt, logPath string) bool {
	args := buildCodexArgs(*session, prompt)
	cmd := exec.Command("codex", args...)
	cmd.Dir = o.root
	output, err := cmd.CombinedOutput()
	if writeErr := os.WriteFile(logPath, output, 0o644); writeErr != nil {
		log.Printf("[%s worker %d] failed writing log: %v", kind, workerID, writeErr)
	}
	if newID := extractSessionID(output); newID != "" && newID != *session {
		*session = newID
		log.Printf("[%s worker %d] captured session %s", kind, workerID, newID)
	}
	if err != nil {
		log.Printf("[%s worker %d] command error: %v (log %s)", kind, workerID, err, logPath)
		return false
	}
	log.Printf("[%s worker %d] completed task (log %s)", kind, workerID, logPath)
	return true
}

type workerPool[T any] struct {
	name    string
	max     int
	runner  func(*worker[T], T)
	idle    chan *worker[T]
	mu      sync.Mutex
	workers []*worker[T]
}

type worker[T any] struct {
	id        int
	pool      *workerPool[T]
	taskCh    chan T
	sessionID string
}

func newWorkerPool[T any](name string, max int, runner func(*worker[T], T)) *workerPool[T] {
	if max < 1 {
		max = 1
	}
	return &workerPool[T]{
		name:   name,
		max:    max,
		runner: runner,
		idle:   make(chan *worker[T], max),
	}
}

func (p *workerPool[T]) Submit(task T) {
	worker := p.acquireWorker()
	worker.taskCh <- task
}

func (p *workerPool[T]) acquireWorker() *worker[T] {
	select {
	case w := <-p.idle:
		return w
	default:
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.workers) < p.max {
		id := len(p.workers) + 1
		w := &worker[T]{
			id:     id,
			pool:   p,
			taskCh: make(chan T),
		}
		p.workers = append(p.workers, w)
		go w.loop()
		return w
	}

	return <-p.idle
}

func (w *worker[T]) loop() {
	for task := range w.taskCh {
		w.pool.runner(w, task)
		w.pool.idle <- w
	}
}

func parseProblemFile(path string) (contest, index string, ok bool) {
	base := filepath.Base(path)
	if !strings.HasPrefix(base, "problem") || !strings.HasSuffix(base, ".txt") {
		return "", "", false
	}
	index = strings.TrimSuffix(strings.TrimPrefix(base, "problem"), ".txt")
	if index == "" {
		return "", "", false
	}
	dir := filepath.Dir(path)
	contest = filepath.Base(dir)
	if contest == "." || contest == "/" || contest == "" {
		return "", "", false
	}
	return contest, index, true
}

func parseSolutionFile(path string) (contest, index string, ok bool) {
	base := filepath.Base(path)
	if !strings.HasSuffix(base, ".go") || strings.HasPrefix(base, "verifier") {
		return "", "", false
	}
	dir := filepath.Dir(path)
	contest = filepath.Base(dir)
	if contest == "." || contest == "/" || contest == "" {
		return "", "", false
	}
	name := strings.TrimSuffix(base, ".go")
	if !strings.HasPrefix(name, contest) {
		return "", "", false
	}
	index = strings.TrimPrefix(name, contest)
	if index == "" {
		return "", "", false
	}
	return contest, index, true
}

func (o *orchestrator) relative(path string) string {
	rel, err := filepath.Rel(o.root, path)
	if err != nil {
		return path
	}
	if rel == "." {
		return rel
	}
	return filepath.ToSlash(rel)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func buildCodexArgs(sessionID, prompt string) []string {
	args := []string{"exec"}
	if sessionID == "" {
		args = append(args, "--full-auto", prompt)
		return args
	}
	args = append(args, "resume", sessionID, prompt)
	return args
}

var sessionPattern = regexp.MustCompile(`session id:\s*([a-f0-9-]+)`)

func extractSessionID(output []byte) string {
	lower := bytes.ToLower(output)
	matches := sessionPattern.FindSubmatch(lower)
	if len(matches) < 2 {
		return ""
	}
	return string(matches[1])
}
