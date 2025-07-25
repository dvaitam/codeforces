package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type testCaseD struct {
	input    string
	expected []int
}

func computeD(n int, arr []int, queries [][2]int) []int {
	res := make([]int, len(queries))
	for qi, q := range queries {
		l, r := q[0]-1, q[1]-1
		last := make(map[int]int)
		best := n + 1
		for i := l; i <= r; i++ {
			v := arr[i]
			if p, ok := last[v]; ok {
				if i-p < best {
					best = i - p
				}
			}
			last[v] = i
		}
		if best == n+1 {
			res[qi] = -1
		} else {
			res[qi] = best
		}
	}
	return res
}

func generateCaseD() testCaseD {
	n := rand.Intn(10) + 1
	m := rand.Intn(10) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(7) - 3 // values -3..3
	}
	queries := make([][2]int, m)
	for i := 0; i < m; i++ {
		l := rand.Intn(n) + 1
		r := rand.Intn(n-l+1) + l
		queries[i] = [2]int{l, r}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for _, q := range queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
	}
	return testCaseD{input: sb.String(), expected: computeD(n, arr, queries)}
}

func parseInts(s string) ([]int, error) {
	parts := strings.Fields(s)
	vals := make([]int, len(parts))
	for i, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		vals[i] = v
	}
	return vals, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 100; i++ {
		tc := generateCaseD()
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, tc.input)
			os.Exit(1)
		}
		vals, err := parseInts(out)
		if err != nil || len(vals) != len(tc.expected) {
			fmt.Fprintf(os.Stderr, "case %d: invalid output\ninput:\n%s", i, tc.input)
			os.Exit(1)
		}
		for j, v := range vals {
			if v != tc.expected[j] {
				fmt.Fprintf(os.Stderr, "case %d: expected %v got %v\ninput:\n%s", i, tc.expected, vals, tc.input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
