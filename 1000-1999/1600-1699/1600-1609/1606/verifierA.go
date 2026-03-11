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

func buildOracle() (string, error) {
	refSrc := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSrc == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	oracle := filepath.Join(os.TempDir(), fmt.Sprintf("oracleA_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", oracle, refSrc)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func genCase(r *rand.Rand) string {
	t := r.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := r.Intn(100) + 1
		for j := 0; j < n; j++ {
			if r.Intn(2) == 0 {
				sb.WriteByte('a')
			} else {
				sb.WriteByte('b')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func countAB(s string) int {
	c := 0
	for i := 0; i+1 < len(s); i++ {
		if s[i] == 'a' && s[i+1] == 'b' {
			c++
		}
	}
	return c
}

func countBA(s string) int {
	c := 0
	for i := 0; i+1 < len(s); i++ {
		if s[i] == 'b' && s[i+1] == 'a' {
			c++
		}
	}
	return c
}

func minSteps(original, modified string) int {
	if len(original) != len(modified) {
		return -1
	}
	c := 0
	for i := range original {
		if original[i] != modified[i] {
			c++
		}
	}
	return c
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		// Parse test cases from input
		lines := strings.Split(strings.TrimSpace(input), "\n")
		t := 0
		fmt.Sscanf(lines[0], "%d", &t)
		originals := lines[1 : t+1]

		// Get oracle output to determine minimum steps
		oracleOut, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		oracleLines := strings.Split(oracleOut, "\n")

		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		gotLines := strings.Split(got, "\n")

		if len(gotLines) != t {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d lines got %d\ninput:\n%s", i, t, len(gotLines), input)
			os.Exit(1)
		}

		for j := 0; j < t; j++ {
			orig := originals[j]
			oracleLine := oracleLines[j]
			gotLine := gotLines[j]

			// Validate: same length as original
			if len(gotLine) != len(orig) {
				fmt.Fprintf(os.Stderr, "case %d test %d: output length %d != input length %d\ninput:\n%s", i, j+1, len(gotLine), len(orig), input)
				os.Exit(1)
			}

			// Validate: only a/b characters
			for _, c := range gotLine {
				if c != 'a' && c != 'b' {
					fmt.Fprintf(os.Stderr, "case %d test %d: invalid character in output\ninput:\n%s", i, j+1, input)
					os.Exit(1)
				}
			}

			// Validate: AB(s) == BA(s)
			if countAB(gotLine) != countBA(gotLine) {
				fmt.Fprintf(os.Stderr, "case %d test %d: AB=%d != BA=%d for output %s\ninput:\n%s", i, j+1, countAB(gotLine), countBA(gotLine), gotLine, input)
				os.Exit(1)
			}

			// Validate: minimum steps
			oracleSteps := minSteps(orig, oracleLine)
			gotSteps := minSteps(orig, gotLine)
			if gotSteps != oracleSteps {
				fmt.Fprintf(os.Stderr, "case %d test %d: got %d steps, expected %d steps\ninput:\n%s", i, j+1, gotSteps, oracleSteps, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All 100 tests passed")
}
