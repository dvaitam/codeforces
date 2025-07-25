package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	n int
	s string
}

func run(bin, input string) (string, error) {
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

func expected(tc testCase) string {
	res := []int{}
	cnt := 0
	for i := 0; i < tc.n; i++ {
		if tc.s[i] == 'B' {
			cnt++
		} else if cnt > 0 {
			res = append(res, cnt)
			cnt = 0
		}
	}
	if cnt > 0 {
		res = append(res, cnt)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprint(len(res)))
	sb.WriteByte('\n')
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	return strings.TrimSpace(sb.String())
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("%d\n%s\n", tc.n, tc.s)
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(100) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = 'B'
		} else {
			b[i] = 'W'
		}
	}
	return testCase{n: n, s: string(b)}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		{1, "B"},
		{1, "W"},
		{3, "BBB"},
		{4, "BWBW"},
	}
	for len(cases) < 105 {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		input := buildInput(tc)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		exp := expected(tc)
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
