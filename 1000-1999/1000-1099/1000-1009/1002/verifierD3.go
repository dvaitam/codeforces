package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const harnessSource = `package main

import (
	"fmt"
	"os"
)

type Qubit int

var qubitState = make([]int, 4)
var initial = make([]int, 3)

func resetInputs(bits int, yInit int) ([]Qubit, Qubit) {
	for i := range qubitState {
		qubitState[i] = 0
	}
	for i := 0; i < 3; i++ {
		val := (bits >> uint(i)) & 1
		qubitState[i] = val
		initial[i] = val
	}
	qubitState[3] = yInit
	return []Qubit{Qubit(0), Qubit(1), Qubit(2)}, Qubit(3)
}

func ensureInputsUnchanged() {
	for i := 0; i < 3; i++ {
		if qubitState[i] != initial[i] {
			fmt.Fprintf(os.Stderr, "input qubit %d changed: expected %d got %d\n", i, initial[i], qubitState[i])
			os.Exit(1)
		}
	}
}

func CCNOT(a, b, target Qubit) {
	if qubitState[int(a)] == 1 && qubitState[int(b)] == 1 {
		qubitState[int(target)] ^= 1
	}
}

func main() {
	results := make([]int, 0, 16)
	for bits := 0; bits < 8; bits++ {
		for y := 0; y <= 1; y++ {
			xs, yq := resetInputs(bits, y)
			MajorityOracle(xs, yq)
			ensureInputsUnchanged()
			results = append(results, qubitState[3])
		}
	}
	for i, v := range results {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(v)
	}
	fmt.Println()
}
`

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD3.go /path/to/solution.go")
		os.Exit(1)
	}
	candidatePath := os.Args[1]
	if filepath.Ext(candidatePath) != ".go" {
		fmt.Fprintln(os.Stderr, "candidate path must point to a Go source file (.go)")
		os.Exit(1)
	}

	refPath, err := findReferenceSolution()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to locate reference solution: %v\n", err)
		os.Exit(1)
	}

	refResults, err := runHarness(refPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference execution failed: %v\n", err)
		os.Exit(1)
	}

	candResults, err := runHarness(candidatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate execution failed: %v\n", err)
		os.Exit(1)
	}

	if len(refResults) != len(candResults) {
		fmt.Fprintf(os.Stderr, "mismatched output lengths: expected %d values got %d\n", len(refResults), len(candResults))
		os.Exit(1)
	}
	for i := range refResults {
		if refResults[i] != candResults[i] {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", i+1, refResults[i], candResults[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(refResults))
}

func findReferenceSolution() (string, error) {
	candidates := []string{
		"1002D3.go",
		filepath.Join(filepath.Dir(os.Args[0]), "1002D3.go"),
	}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "1002D3.go"))
	}
	for _, path := range candidates {
		if path == "" {
			continue
		}
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("1002D3.go not found in expected locations")
}

func runHarness(sourcePath string) ([]int, error) {
	tmpDir, err := os.MkdirTemp("", "verifier-d3-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	harnessFile := filepath.Join(tmpDir, "harness.go")
	if err := os.WriteFile(harnessFile, []byte(harnessSource), 0644); err != nil {
		return nil, fmt.Errorf("failed to write harness: %v", err)
	}

	code, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %v", sourcePath, err)
	}
	solutionFile := filepath.Join(tmpDir, "solution.go")
	if err := os.WriteFile(solutionFile, code, 0644); err != nil {
		return nil, fmt.Errorf("failed to copy solution: %v", err)
	}

	cmd := exec.Command("go", "run", harnessFile, solutionFile)
	cmd.Dir = tmpDir
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("go run failed: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}

	fields := strings.Fields(stdout.String())
	if len(fields) == 0 {
		return nil, fmt.Errorf("no output from harness")
	}
	results := make([]int, len(fields))
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("failed to parse output %q: %v\nfull output:\n%s", f, err, stdout.String())
		}
		results[i] = val
	}
	return results, nil
}
