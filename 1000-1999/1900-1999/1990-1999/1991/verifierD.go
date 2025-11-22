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
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", oracle, "1991D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	return fmt.Sprintf("1\n%d\n", n)
}

func isPrime(x int) bool {
	if x < 2 {
		return false
	}
	if x == 2 || x == 3 {
		return true
	}
	if x%2 == 0 {
		return false
	}
	for i := 3; i*i <= x; i += 2 {
		if x%i == 0 {
			return false
		}
	}
	return true
}

func parseOutput(out string, n int) (int, []int, error) {
	var k int
	colors := make([]int, n)
	reader := strings.NewReader(out)
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return 0, nil, fmt.Errorf("failed to read k: %w", err)
	}
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(reader, &colors[i]); err != nil {
			return 0, nil, fmt.Errorf("failed to read color %d: %w", i+1, err)
		}
	}
	for i, c := range colors {
		if c < 1 || c > k {
			return 0, nil, fmt.Errorf("color %d out of range: %d", i+1, c)
		}
	}
	return k, colors, nil
}

func validColoring(colors []int) bool {
	n := len(colors)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if colors[i] == colors[j] && isPrime((i+1)^(j+1)) {
				return false
			}
		}
	}
	return true
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, errOut.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
	cases := make([]string, 0, 100)
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}
	for i, in := range cases {
		exp, err := run(oracle, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		var n int
		fmt.Sscanf(in, "%*d\n%d", &n)
		expK, _, err := parseOutput(exp, n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle parse error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		gotK, colors, err := parseOutput(got, n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid output on case %d: %v\ninput:\n%soutput:\n%s", i+1, err, in, got)
			os.Exit(1)
		}
		if gotK != expK {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected k: %d\ngot k: %d\ninput:\n%s", i+1, expK, gotK, in)
			os.Exit(1)
		}
		if !validColoring(colors) {
			fmt.Fprintf(os.Stderr, "case %d failed: coloring invalid\ninput:\n%soutput:\n%s", i+1, in, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
