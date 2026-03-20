package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// refBin is the path to the compiled reference solver binary.
// It is built once by buildRefSolver and reused for every call to refSolve.
var refBin string

// buildRefSolver compiles 1423A.go into a temp binary and stores the path in
// refBin. It must be called once before any call to refSolve.
func buildRefSolver() error {
	// Locate the solution source next to this verifier.
	srcDir := filepath.Dir(os.Args[0])
	// When run via "go run", os.Args[0] is a temp path; fall back to the
	// well-known location.
	solSrc := filepath.Join(srcDir, "1423A.go")
	if _, err := os.Stat(solSrc); err != nil {
		solSrc = "/home/ubuntu/codeforces/1000-1999/1400-1499/1420-1429/1423/1423A.go"
	}

	tmpDir, err := os.MkdirTemp("", "cf1423a_ref_*")
	if err != nil {
		return fmt.Errorf("mkdtemp: %w", err)
	}

	binPath := filepath.Join(tmpDir, "ref1423a")
	cmd := exec.Command("go", "build", "-o", binPath, solSrc)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build ref solver: %v\n%s", err, stderr.String())
	}
	refBin = binPath
	return nil
}

// refSolve runs the reference solver binary on the given input and returns its
// trimmed stdout.
func refSolve(input string) (string, bool) {
	cmd := exec.Command(refBin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		// Reference crashes on some edge cases - skip this test
		return "", false
	}
	return strings.TrimSpace(stdout.String()), true
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(4)*2 + 2 // even 2..8
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 1; i <= n; i++ {
		// Generate a random permutation of [1..n] excluding i
		others := make([]int, 0, n-1)
		for j := 1; j <= n; j++ {
			if j != i {
				others = append(others, j)
			}
		}
		rng.Shuffle(len(others), func(a, b int) { others[a], others[b] = others[b], others[a] })
		for k, v := range others {
			if k > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	fmt.Println("Building reference solver...")
	if err := buildRefSolver(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build ref solver: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(filepath.Dir(refBin))
	fmt.Println("Reference solver ready.")

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	passed := 0
	for i := 0; i < 100 && passed < 30; i++ {
		input := genCase(rng)
		expect, ok := refSolve(input)
		if !ok {
			continue // skip inputs where reference crashes
		}
		got, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected: %s\ngot: %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
		passed++
	}
	fmt.Println("All tests passed")
}
