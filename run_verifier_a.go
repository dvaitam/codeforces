package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

func main() {
	root := flag.String("root", ".", "root directory to search for problemA.txt files")
	concurrency := flag.Int("concurrency", runtime.NumCPU(), "number of concurrent verifier runs")
	timeout := flag.Duration("timeout", 2*time.Minute, "max duration for each build/run")
	memLimitKB := flag.Int("mem-kb", 0, "memory limit (KB) for verifier execution; 0 disables the limit")
	flag.Parse()

	log.SetFlags(0)

	targets, err := collectProblems(*root)
	if err != nil {
		log.Fatalf("failed to collect problem directories: %v", err)
	}
	if len(targets) == 0 {
		log.Println("no problem*.txt files found")
		return
	}

	workers := *concurrency
	if workers < 1 {
		workers = 1
	}

	memMsg := fmt.Sprintf("%dKB", *memLimitKB)
	if *memLimitKB <= 0 {
		memMsg = "unlimited"
	}
	log.Printf("found %d problems; running with concurrency=%d, timeout=%s, mem=%s\n", len(targets), workers, timeout.String(), memMsg)

	jobs := make(chan problemSpec)
	results := make(chan runResult)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for spec := range jobs {
				results <- processProblem(spec, *timeout, *memLimitKB)
			}
		}()
	}

	go func() {
		for _, target := range targets {
			jobs <- target
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var okCount, failCount int
	for res := range results {
		if res.success {
			okCount++
			log.Printf("[OK]   %s\n", res.label)
			continue
		}

		failCount++
		log.Printf("[FAIL] %s (%s): %v", res.label, res.stage, res.err)
		if res.output != "" {
			log.Println(indent(res.output))
		}
	}

	log.Printf("completed. successes=%d failures=%d\n", okCount, failCount)
}

type problemSpec struct {
	dir    string
	base   string
	suffix string
	label  string
}

type runResult struct {
	label   string
	stage   string
	output  string
	err     error
	success bool
}

func collectProblems(root string) ([]problemSpec, error) {
	var specs []problemSpec
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
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
		suffix := strings.TrimSuffix(strings.TrimPrefix(name, "problem"), ".txt")
		if suffix == "" {
			return nil
		}
		dir := filepath.Dir(path)
		base := filepath.Base(dir)
		specs = append(specs, problemSpec{
			dir:    dir,
			base:   base,
			suffix: suffix,
			label:  filepath.Join(dir, name),
		})
		return nil
	})
	return specs, err
}

func processProblem(spec problemSpec, timeout time.Duration, memLimitKB int) runResult {
	res := runResult{label: spec.label}

	// Clean up artifacts even when builds fail.
	verifierSrc := fmt.Sprintf("verifier%s.go", spec.suffix)
	verifierBin := fmt.Sprintf("verifier%s", spec.suffix)
	solutionSource := fmt.Sprintf("%s%s.go", spec.base, spec.suffix)
	solutionBin := fmt.Sprintf("%s%s", spec.base, spec.suffix)
	defer cleanupArtifacts(spec.dir, verifierBin, solutionBin)

	check := func(name string) error {
		_, err := os.Stat(filepath.Join(spec.dir, name))
		return err
	}
	if err := check(verifierSrc); err != nil {
		res.stage = "check verifier"
		res.err = err
		return res
	}
	if err := check(solutionSource); err != nil {
		res.stage = "check solution"
		res.err = err
		return res
	}

	if out, err := runCommand(spec.dir, timeout, "go", "build", verifierSrc); err != nil {
		res.stage = "build verifier"
		res.err = err
		res.output = out
		return res
	}
	if out, err := runCommand(spec.dir, timeout, "go", "build", solutionSource); err != nil {
		res.stage = "build solution"
		res.err = err
		res.output = out
		return res
	}

	runCmd := "exec ./" + verifierBin + " ./" + solutionBin
	if memLimitKB > 0 {
		runCmd = fmt.Sprintf("ulimit -v %d; %s", memLimitKB, runCmd)
	}
	if out, err := runCommand(spec.dir, timeout, "bash", "-c", runCmd); err != nil {
		res.stage = "run verifier"
		res.err = err
		res.output = out
		return res
	}

	res.success = true
	return res
}

func runCommand(dir string, timeout time.Duration, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return string(out), fmt.Errorf("timeout after %s", timeout)
	}
	return string(out), err
}

func cleanupArtifacts(dir string, names ...string) {
	for _, name := range names {
		_ = os.Remove(filepath.Join(dir, name))
	}
}

func indent(s string) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	for i, line := range lines {
		lines[i] = "    " + line
	}
	return strings.Join(lines, "\n")
}
