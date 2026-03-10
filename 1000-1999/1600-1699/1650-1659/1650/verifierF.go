package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func buildExecutable(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "bin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), path)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, stderr.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func oracle(exe string, test string) (string, error) {
	return run(exe, "1\n"+test)
}

// validateAnswer checks if a candidate answer is valid for the given test case.
// Returns true if the answer is correct.
func validateAnswer(tc string, got string) bool {
	lines := strings.Split(strings.TrimSpace(tc), "\n")
	parts := strings.Fields(lines[0])
	n, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	a := make([]int64, n)
	aFields := strings.Fields(lines[1])
	for i := 0; i < n; i++ {
		v, _ := strconv.ParseInt(aFields[i], 10, 64)
		a[i] = v
	}
	type option struct {
		e int
		t int64
		p int64
	}
	opts := make([]option, m)
	for i := 0; i < m; i++ {
		fields := strings.Fields(lines[2+i])
		e, _ := strconv.Atoi(fields[0])
		t, _ := strconv.ParseInt(fields[1], 10, 64)
		p, _ := strconv.ParseInt(fields[2], 10, 64)
		opts[i] = option{e: e - 1, t: t, p: p}
	}

	gotLines := strings.Split(strings.TrimSpace(got), "\n")
	if len(gotLines) == 0 {
		return false
	}
	k, err := strconv.Atoi(strings.TrimSpace(gotLines[0]))
	if err != nil {
		return false
	}
	if k == -1 {
		// Check if -1 is actually correct (no valid solution exists)
		// We trust the oracle for this
		return false
	}
	if k == 0 {
		// Using 0 options: check all tasks already at >= 100% (they start at 0, so impossible unless n==0)
		// With 0 options, progress for all tasks is 0, which is < 100
		if n == 0 {
			return true
		}
		return false
	}
	if len(gotLines) < 2 {
		return false
	}
	ids := strings.Fields(gotLines[1])
	if len(ids) != k {
		return false
	}
	// Check all ids are valid and unique
	used := make(map[int]bool)
	chosenOpts := make([]option, 0, k)
	chosenIds := make([]int, 0, k)
	for _, idStr := range ids {
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 || id > m {
			return false
		}
		if used[id] {
			return false
		}
		used[id] = true
		chosenOpts = append(chosenOpts, opts[id-1])
		chosenIds = append(chosenIds, id)
	}
	// Simulate: execute options in the given order, track time and progress
	progress := make([]int64, n)
	var curTime int64
	for _, opt := range chosenOpts {
		curTime += opt.t
		progress[opt.e] += opt.p
	}
	// Check each task: must have >= 100% progress by its deadline
	var timeUsed int64
	// Actually, we need to check that for each task i, the total time of all options
	// chosen for tasks 1..i is <= a[i]. Tasks must be completed in order.
	// Re-simulate respecting deadlines
	// Sort chosen options by task order (tasks with earlier deadlines first)
	// Actually, the options should be executed in the given order, and we need
	// to verify that cumulative time at completion of each task doesn't exceed its deadline.

	// Group options by task
	taskProgress := make([]int64, n)
	taskTime := make([]int64, n)
	for _, opt := range chosenOpts {
		taskProgress[opt.e] += opt.p
		taskTime[opt.e] += opt.t
	}

	// Check all tasks have >= 100% progress
	for i := 0; i < n; i++ {
		if taskProgress[i] < 100 {
			return false
		}
	}

	// Check timing: cumulative time up to task i must be <= a[i]
	timeUsed = 0
	for i := 0; i < n; i++ {
		timeUsed += taskTime[i]
		if timeUsed > a[i] {
			return false
		}
	}

	return true
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + n
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d", rng.Intn(50)+1)
		if i+1 < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		e := rng.Intn(n) + 1
		t := rng.Intn(20) + 1
		p := rng.Intn(100) + 1
		fmt.Fprintf(&sb, "%d %d %d\n", e, t, p)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	binPath := os.Args[1]
	bin, cleanup, err := buildExecutable(binPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare binary: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()
	oracleSrc := os.Getenv("REFERENCE_SOURCE_PATH")
	if oracleSrc == "" {
		log.Fatal("REFERENCE_SOURCE_PATH environment variable is not set")
	}
	oracleExe, oracleCleanup, err := buildExecutable(oracleSrc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer oracleCleanup()
	rng := rand.New(rand.NewSource(6))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		expected, err := oracle(oracleExe, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed: %v\n", err)
			os.Exit(1)
		}
		got, err := run(bin, "1\n"+tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if got != expected {
			// The problem allows multiple valid answers, so validate the candidate's answer
			expectedIsNeg := strings.TrimSpace(expected) == "-1"
			gotIsNeg := strings.TrimSpace(got) == "-1"
			if expectedIsNeg && !gotIsNeg {
				// Oracle says impossible but candidate found a solution - check it
				if !validateAnswer(tc, got) {
					fmt.Fprintf(os.Stderr, "case %d failed: oracle says -1, candidate answer invalid\ninput:\n%s", i+1, tc)
					os.Exit(1)
				}
			} else if !expectedIsNeg && gotIsNeg {
				fmt.Fprintf(os.Stderr, "case %d failed: expected valid answer got -1\ninput:\n%s", i+1, tc)
				os.Exit(1)
			} else if !expectedIsNeg && !gotIsNeg {
				if !validateAnswer(tc, got) {
					fmt.Fprintf(os.Stderr, "case %d failed: candidate answer invalid\ninput:\n%sgot:\n%s\n", i+1, tc, got)
					os.Exit(1)
				}
			}
			// Both -1 means both agree no solution exists
		}
	}
	fmt.Println("All tests passed")
}
