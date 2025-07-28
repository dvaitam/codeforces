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

type testCase struct {
	n       int
	queries []string
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "1638E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func generateCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 100)
	for i := range cases {
		n := rng.Intn(20) + 1
		qnum := rng.Intn(30) + 1
		qs := make([]string, qnum)
		for j := 0; j < qnum; j++ {
			t := rng.Intn(3)
			switch t {
			case 0:
				l := rng.Intn(n) + 1
				r := rng.Intn(n) + 1
				c := rng.Intn(n) + 1
				qs[j] = fmt.Sprintf("C %d %d %d", l, r, c)
			case 1:
				c := rng.Intn(n) + 1
				x := rng.Intn(1000) + 1
				qs[j] = fmt.Sprintf("A %d %d", c, x)
			default:
				idx := rng.Intn(n) + 1
				qs[j] = fmt.Sprintf("Q %d", idx)
			}
		}
		cases[i] = testCase{n: n, queries: qs}
	}
	return cases
}

func runCase(oracle, bin string, tc testCase) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.queries)))
	for _, q := range tc.queries {
		sb.WriteString(q)
		sb.WriteByte('\n')
	}
	input := sb.String()

	cmdO := exec.Command(oracle)
	cmdO.Stdin = strings.NewReader(input)
	var outO bytes.Buffer
	cmdO.Stdout = &outO
	if err := cmdO.Run(); err != nil {
		return fmt.Errorf("oracle runtime error: %v", err)
	}
	expected := strings.TrimSpace(outO.String())

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	cases := generateCases()
	for i, c := range cases {
		if err := runCase(oracle, bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
