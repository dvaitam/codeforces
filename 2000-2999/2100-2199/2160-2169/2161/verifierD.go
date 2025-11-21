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

const refPath = "2000-2999/2100-2199/2160-2169/2161/2161D.go"

type testCase struct {
	n   int
	arr []int
}

func genCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	totalN := 0
	for len(cases) < 200 && totalN < 4000 {
		n := rng.Intn(30) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(n) + 1
		}
		cases = append(cases, testCase{n: n, arr: arr})
		totalN += n
	}
	if len(cases) == 0 {
		cases = append(cases, testCase{n: 1, arr: []int{1}})
	}
	return cases
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cases))
	for _, tc := range cases {
		fmt.Fprintln(&sb, tc.n)
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runProgram(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/2161D_binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	cases := genCases()
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

	expTokens := strings.Fields(expected)
	actTokens := strings.Fields(actual)
	if len(actTokens) < len(expTokens) {
		fmt.Fprintf(os.Stderr, "not enough outputs: got %d expected %d\n", len(actTokens), len(expTokens))
		os.Exit(1)
	}
	for i := range expTokens {
		expVal, err1 := strconv.ParseInt(expTokens[i], 10, 64)
		actVal, err2 := strconv.ParseInt(actTokens[i], 10, 64)
		if err1 != nil || err2 != nil || expVal != actVal {
			fmt.Fprintf(os.Stderr, "mismatch at test %d: expected %s got %s\nn=%d arr=%v\n", i+1, expTokens[i], actTokens[i], cases[i].n, cases[i].arr)
			os.Exit(1)
		}
	}
	if len(actTokens) != len(expTokens) {
		fmt.Fprintf(os.Stderr, "extra output tokens detected\n")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
