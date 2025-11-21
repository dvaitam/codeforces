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

type ring struct {
	x, y int
	r, R int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refSrc, err := locateReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refBin, err := buildReference(refSrc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, input := range tests {
		want, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if normalize(got) != normalize(want) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d\nInput:\n%sExpected:\n%sGot:\n%s\n", idx+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"199B.go",
		filepath.Join("0-999", "100-199", "190-199", "199", "199B.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 199B.go relative to working directory")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref199B_%d.bin", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", outPath, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func normalize(out string) string {
	return strings.TrimSpace(out)
}

func generateTests() []string {
	var tests []string
	tests = append(tests, formatInput(
		ring{0, 0, 1, 2},
		ring{5, 0, 1, 3},
	))
	tests = append(tests, formatInput(
		ring{0, 0, 1, 5},
		ring{3, 4, 2, 6},
	))
	tests = append(tests, formatInput(
		ring{-10, -10, 1, 10},
		ring{10, 10, 5, 20},
	))
	tests = append(tests, formatInput(
		ring{0, 0, 10, 20},
		ring{21, 0, 1, 2},
	))
	tests = append(tests, formatInput(
		ring{0, 0, 1, 50},
		ring{50, 50, 49, 50},
	))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 200 {
		r1 := randRing(rng)
		r2 := randRing(rng)
		for r1.x == r2.x && r1.y == r2.y {
			r2 = randRing(rng)
		}
		tests = append(tests, formatInput(r1, r2))
	}
	return tests
}

func randRing(rng *rand.Rand) ring {
	r := rng.Intn(100) + 1
	R := rng.Intn(100-r) + r + 1
	return ring{
		x: rng.Intn(201) - 100,
		y: rng.Intn(201) - 100,
		r: r,
		R: R,
	}
}

func formatInput(a, b ring) string {
	return fmt.Sprintf("%d %d %d %d\n%d %d %d %d\n", a.x, a.y, a.r, a.R, b.x, b.y, b.r, b.R)
}
