package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const refPath = "2000-2999/2100-2199/2110-2119/2113/2113A.go"

type testCase struct {
	k, a, b, x, y int64
}

func genCase(rng *rand.Rand) testCase {
	k := rng.Int63n(1_000_000_000) + 1
	a := rng.Int63n(1_000_000_000) + 1
	b := rng.Int63n(1_000_000_000) + 1
	x := rng.Int63n(1_000_000_000) + 1
	y := rng.Int63n(1_000_000_000) + 1
	return testCase{k: k, a: a, b: b, x: x, y: y}
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&sb, "%d %d %d %d %d\n", tc.k, tc.a, tc.b, tc.x, tc.y)
	}
	return sb.String()
}

func runProgram(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", path)
	} else {
		cmd = exec.CommandContext(ctx, path)
	}
	cmd.Stdin = strings.NewReader(input)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout running %s", path)
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(outBuf.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/2113A_binary")
		os.Exit(1)
	}
	cand := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	for i := 0; i < 200; i++ {
		cases = append(cases, genCase(rng))
	}

	input := buildInput(cases)

	refAbs, _ := filepath.Abs(refPath)
	candAbs := cand
	if !filepath.IsAbs(candAbs) {
		candAbs, _ = filepath.Abs(cand)
	}

	expected, err := runProgram(refAbs, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run reference: %v\n", err)
		os.Exit(1)
	}
	actual, err := runProgram(candAbs, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}

	expLines := strings.Fields(expected)
	actLines := strings.Fields(actual)
	if len(actLines) < len(expLines) {
		fmt.Fprintf(os.Stderr, "not enough outputs: got %d expected %d\n", len(actLines), len(expLines))
		os.Exit(1)
	}

	for i := range expLines {
		expVal, err1 := strconv.ParseInt(expLines[i], 10, 64)
		actVal, err2 := strconv.ParseInt(actLines[i], 10, 64)
		if err1 != nil || err2 != nil || expVal != actVal {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %s got %s\nparams: %v\n", i+1, expLines[i], actLines[i], cases[i])
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed")
}
