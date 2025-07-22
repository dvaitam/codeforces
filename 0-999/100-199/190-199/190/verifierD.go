package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func solve(n, k int, arr []int) string {
	pos := make(map[int][]int)
	curMax := 0
	var ans int64
	for i, v := range arr {
		idx := i + 1
		lst := pos[v]
		lst = append(lst, idx)
		pos[v] = lst
		if len(lst) >= k {
			t := lst[len(lst)-k]
			if t > curMax {
				curMax = t
			}
		}
		if curMax > 0 {
			ans += int64(curMax)
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generateCases() []testCase {
	rand.Seed(4)
	cases := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(20) + 1
		k := rand.Intn(n) + 1
		arr := make([]int, n)
		buf := bytes.Buffer{}
		fmt.Fprintf(&buf, "%d %d\n", n, k)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(10)
			if j > 0 {
				fmt.Fprint(&buf, " ")
			}
			fmt.Fprint(&buf, arr[j])
		}
		buf.WriteByte('\n')
		cases[i] = testCase{input: buf.String(), expected: solve(n, k, arr)}
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierD.go <binary>")
		os.Exit(1)
	}
	cases := generateCases()
	for i, tc := range cases {
		out, err := runBinary(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%s\nexpected:%s\nactual:%s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
