package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func buildOracle() (string, error) {
	dir, _ := os.Getwd()
	oracle := filepath.Join(dir, "oracleF")
	if out, err := exec.Command("go", "build", "-o", oracle, "1359F.go").CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	return oracle, nil
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type testCase struct {
	n    int
	vals []string
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	vals := make([]string, n)
	for i := 0; i < n; i++ {
		x := rng.Intn(21) - 10
		y := rng.Intn(21) - 10
		dx := rng.Intn(10) + 1
		if rng.Intn(2) == 0 {
			dx = -dx
		}
		dy := rng.Intn(10) + 1
		if rng.Intn(2) == 0 {
			dy = -dy
		}
		s := rng.Intn(10) + 1
		vals[i] = fmt.Sprintf("%d %d %d %d %d", x, y, dx, dy, s)
	}
	return testCase{n, vals}
}

func compareFloats(a, b float64) bool {
	diff := math.Abs(a - b)
	denom := math.Max(1.0, math.Abs(a))
	return diff/denom <= 1e-6
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for _, v := range tc.vals {
			sb.WriteString(v + "\n")
		}
		input := sb.String()
		expStr, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failure on case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		gotStr, err := runProg(exe, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if expStr == "No show :(" {
			if gotStr != expStr {
				fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:%s", i+1, expStr, gotStr, input)
				os.Exit(1)
			}
			continue
		}
		expF, _ := strconv.ParseFloat(expStr, 64)
		gotF, err := strconv.ParseFloat(gotStr, 64)
		if err != nil || !compareFloats(expF, gotF) {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:%s", i+1, expStr, gotStr, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
