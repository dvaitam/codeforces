package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refPath = "2000-2999/2100-2199/2110-2119/2116/2116A.go"

type testCase struct {
	a, b, c, d int64
}

func genCase(rng *rand.Rand) testCase {
	return testCase{
		a: rng.Int63n(50) + 1,
		b: rng.Int63n(50) + 1,
		c: rng.Int63n(50) + 1,
		d: rng.Int63n(50) + 1,
	}
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&sb, "%d %d %d %d\n", tc.a, tc.b, tc.c, tc.d)
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
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout running %s", path)
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/2116A_binary")
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
		if expLines[i] != actLines[i] {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %s got %s (a=%d b=%d c=%d d=%d)\n",
				i+1, expLines[i], actLines[i], cases[i].a, cases[i].b, cases[i].c, cases[i].d)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed")
}
